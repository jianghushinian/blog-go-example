// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// A Pool is a set of temporary objects that may be individually saved and
// retrieved.
//
// Any item stored in the Pool may be removed automatically at any time without
// notification. If the Pool holds the only reference when this happens, the
// item might be deallocated.
//
// A Pool is safe for use by multiple goroutines simultaneously.
//
// Pool's purpose is to cache allocated but unused items for later reuse,
// relieving pressure on the garbage collector. That is, it makes it easy to
// build efficient, thread-safe free lists. However, it is not suitable for all
// free lists.
//
// An appropriate use of a Pool is to manage a group of temporary items
// silently shared among and potentially reused by concurrent independent
// clients of a package. Pool provides a way to amortize allocation overhead
// across many clients.
//
// An example of good use of a Pool is in the fmt package, which maintains a
// dynamically-sized store of temporary output buffers. The store scales under
// load (when many goroutines are actively printing) and shrinks when
// quiescent.
//
// On the other hand, a free list maintained as part of a short-lived object is
// not a suitable use for a Pool, since the overhead does not amortize well in
// that scenario. It is more efficient to have such objects implement their own
// free list.
//
// A Pool must not be copied after first use.
//
// In the terminology of [the Go memory model], a call to Put(x) “synchronizes before”
// a call to [Pool.Get] returning that same value x.
// Similarly, a call to New returning x “synchronizes before”
// a call to Get returning that same value x.
//
// [the Go memory model]: https://go.dev/ref/mem
type Pool struct {
	// 禁止复制
	noCopy noCopy

	// 空闲对象，poolLocal 指针类型
	local unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	// 数组大小
	localSize uintptr // size of the local array

	// 回收站，poolLocal 指针类型
	victim unsafe.Pointer // local from previous cycle
	// 数组大小
	victimSize uintptr // size of victims array

	// New 是一个可选的函数，调用 Get 方法时，如果缓存池中没有可用对象，则调用此方法生成一个新的值并返回，否则返回 nil。
	// 该函数不能在并发调用 Get 时被修改。
	New func() any
}

// Local per-P Pool appendix.
type poolLocalInternal struct {
	// 私有对象
	private any // Can be used only by the respective P.
	// 共享队列，这是一个 lock-free 双向队列
	shared poolChain // Local P can pushHead/popHead; any P can popTail.
}

type poolLocal struct {
	poolLocalInternal

	// Prevents false sharing on widespread platforms with
	// 128 mod (cache line size) = 0 .
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

// from runtime
//
//go:linkname runtime_randn runtime.randn
func runtime_randn(n uint32) uint32

var poolRaceHash [128]uint64

// poolRaceAddr returns an address to use as the synchronization point
// for race detector logic. We don't use the actual pointer stored in x
// directly, for fear of conflicting with other synchronization on that address.
// Instead, we hash the pointer to get an index into poolRaceHash.
// See discussion on golang.org/cl/31589.
func poolRaceAddr(x any) unsafe.Pointer {
	ptr := uintptr((*[2]unsafe.Pointer)(unsafe.Pointer(&x))[1])
	h := uint32((uint64(uint32(ptr)) * 0x85ebca6b) >> 16)
	return unsafe.Pointer(&poolRaceHash[h%uint32(len(poolRaceHash))])
}

// Put 添加一个元素到池中
func (p *Pool) Put(x any) {
	if x == nil { // x 为 nil 直接返回
		return
	}

	// if race.Enabled {
	// 	if runtime_randn(4) == 0 {
	// 		// Randomly drop x on floor.
	// 		return
	// 	}
	// 	race.ReleaseMerge(poolRaceAddr(x))
	// 	race.Disable()
	// }

	// 把当前 goroutine 固定在当前的 P 上
	// 这样，在当前 goroutine 中操作当前 P 上的任何对象都无需加锁
	// 因为在一个 P 上，同一时刻只会运行一个 goroutine，不会有并发问题
	// 拿到 local 对象（*poolLocal）和当前 P ID
	l, _ := p.pin()

	if l.private == nil {
		l.private = x // 如果 private 为 nil，则直接将 x 赋值给它
	} else {
		l.shared.pushHead(x) // 否则，将 x push 到共享队列队头
	}

	// 将当前 goroutine 从当前 P 上解除固定
	runtime_procUnpin()

	// if race.Enabled {
	// 	race.Enable()
	// }
}

// Get selects an arbitrary item from the [Pool], removes it from the Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to [Pool.Put] and the values returned by Get.
//
// If Get would otherwise return nil and p.New is non-nil, Get returns the result of calling p.New.

// Get 从 [Pool] 中选择一个任意项，将其从 Pool 中移除，然后返回给调用者。
// Get 可以选择忽略池并将其视为空。
// 调用者不应假设传递给 [Pool.Put] 的值与 Get 返回的值之间存在任何关系。
//
// 如果 Get 返回 nil 且 p.New 非零，则 Get 返回调用 p.New 的结果。
func (p *Pool) Get() any {
	// if race.Enabled {
	// 	race.Disable()
	// }

	// 把当前 goroutine 固定在当前的 P 上
	// 拿到 local 对象（*poolLocal，该 P 的本地池）和当前 P ID
	l, pid := p.pin()

	// 空闲对象搜索路径
	// 1. 从 private 中获取对象
	// 2. 从本地共享队列中获取对象
	// 3. 慢路径（则尝试从其他 P 窃取或从 victim 缓存获取）

	// 获取当前 P 中的 private
	x := l.private
	l.private = nil
	if x == nil { // private 不存在
		// 尝试从当前 P 的共享队列中弹出空闲对象
		// 因为 shared 队列只有所属的 P 会操作头部（生产者），所以 popHead 操作也无需加锁
		x, _ = l.shared.popHead()
		if x == nil { // 触发慢路径
			// 当前 P 的本地池为空，则尝试从其他 P 窃取或从 victim 缓存获取
			x = p.getSlow(pid)
		}
	}

	// 解除 pin
	runtime_procUnpin()

	// if race.Enabled {
	// 	race.Enable()
	// 	if x != nil {
	// 		race.Acquire(poolRaceAddr(x))
	// 	}
	// }

	if x == nil && p.New != nil {
		x = p.New() // 如果所有缓存都未找到对象，且用户提供了 New 函数，则创建一个新对象
	}
	return x
}

func (p *Pool) getSlow(pid int) any {
	// See the comment in pin regarding ordering of the loads.
	size := runtime_LoadAcquintptr(&p.localSize) // load-acquire
	locals := p.local                            // load-consume

	// 尝试从其他 P 的共享队列中窃取一个元素
	for i := 0; i < int(size); i++ {
		l := indexLocal(locals, (pid+i+1)%int(size)) // 计算其他 P 的索引
		if x, _ := l.shared.popTail(); x != nil {    // 从尾部窃取
			return x
		}
	}

	// 如果窃取也失败了，则转而检查 victim 缓存

	// Try the victim cache. We do this after attempting to steal
	// from all primary caches because we want objects in the
	// victim cache to age out if at all possible.
	size = atomic.LoadUintptr(&p.victimSize) // 获取 victim 缓存大小
	if uintptr(pid) >= size {
		return nil
	}
	locals = p.victim
	l := indexLocal(locals, pid)
	if x := l.private; x != nil { // 先检查 victim 的 private
		l.private = nil // 从 victim 中移除后再返回
		return x
	}
	for i := 0; i < int(size); i++ { // 再检查其他 P 的 victim 的 shared
		l := indexLocal(locals, (pid+i)%int(size)) // 计算其他 P 的索引
		if x, _ := l.shared.popTail(); x != nil {  // 从尾部窃取
			return x
		}
	}

	// 如果 victim 中也没找到，则返回 nil

	atomic.StoreUintptr(&p.victimSize, 0) // 标记 victim 为空

	return nil
}

// 将当前 goroutine 固定（pin）到其运行的 P（逻辑处理器）上，并返回该 P 对应的本地缓存池 (*poolLocal) 和 P 的 ID
// 调用方必须在使用完成后调用 runtime_procUnpin() 取消固定
func (p *Pool) pin() (*poolLocal, int) {
	// Check whether p is nil to get a panic.
	// Otherwise the nil dereference happens while the m is pinned,
	// causing a fatal error rather than a panic.
	if p == nil { // 空指针检查
		panic("nil Pool")
	}

	// 固定 P，调用 runtime 函数，禁止当前 G 被抢占，并将其固定到当前 P，同时返回 P 的 ID
	// 这是后续无锁操作的基础
	pid := runtime_procPin()
	// In pinSlow we store to local and then to localSize, here we load in opposite order.
	// Since we've disabled preemption, GC cannot happen in between.
	// Thus here we must observe local at least as large localSize.
	// We can observe a newer/larger local, it is fine (we must observe its zero-initialized-ness).
	// 原子加载本地池信息
	s := runtime_LoadAcquintptr(&p.localSize) // load-acquire
	l := p.local                              // load-consume
	if uintptr(pid) < s {                     // 快速路径（常见情况）
		// 如果当前 P 的 ID 在 local 数组的有效大小范围内，则通过 indexLocal 函数计算地址，直接返回对应的 poolLocal 和 pid
		return indexLocal(l, pid), pid
	}

	// 慢路径（初始化或扩容）
	return p.pinSlow()
}

// pin 方法的“慢路径”（slow path）
// 负责在特定情况下初始化或重新分配 Pool 的本地存储数组 (local)，
// 并确保该 Pool 被注册到全局的 allPools 列表中以便垃圾回收 (GC) 时进行清理
func (p *Pool) pinSlow() (*poolLocal, int) {
	// Retry under the mutex.
	// Can not lock the mutex while pinned.
	runtime_procUnpin() // 解除当前 G 与 P 的绑定，为获取全局锁做准备

	allPoolsMu.Lock() // 加全局互斥锁，保护对 allPools 和 Pool 的 local 等字段的并发访问
	defer allPoolsMu.Unlock()

	pid := runtime_procPin() // 重新固定 G 到 P
	// poolCleanup won't be called while we are pinned.
	s := p.localSize
	l := p.local
	if uintptr(pid) < s { // 在锁保护下再次检查 local 数组是否已由其他 goroutine 初始化（双重检查锁定模式）
		return indexLocal(l, pid), pid
	}

	// 如果 Pool 尚未注册，则将其添加到 allPools 全局切片中，以便后续 GC 时能执行 poolCleanup 清理其缓存
	if p.local == nil {
		allPools = append(allPools, p)
	}

	// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal, size) // 根据当前的 GOMAXPROCS（即 P 的数量）创建一个新的 poolLocal 数组
	// 记录初始化的 poolLocal 数组
	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
	runtime_StoreReluintptr(&p.localSize, uintptr(size))     // store-release
	return &local[pid], pid                                  // 返回新创建的、当前 P 对应的 *poolLocal 和 P 的 ID
}

// poolCleanup 本应是一个内部实现细节，但许多广泛使用的包通过 linkname 方式访问了它
// “不光彩名单”中的著名成员包括：
//   - github.com/bytedance/gopkg
//   - github.com/songzhibin97/gkit
//
// 不要移除或更改此函数的类型签名
// 参见 go.dev/issue/67401
//
//go:linkname poolCleanup
func poolCleanup() {
	// 此函数在垃圾回收（GC）开始，程序暂停（STW）时被调用
	// 它自身一定不能分配内存，并且很可能不应调用任何运行时函数（runtime functions）

	// Because the world is stopped, no pool user can be in a pinned section (in effect, this has all Ps pinned).

	// Drop victim caches from all pools.
	for _, p := range oldPools {
		p.victim = nil // 清空回收站
		p.victimSize = 0
	}

	// Move primary cache to victim cache.
	for _, p := range allPools {
		p.victim = p.local // 从主缓存移到回收站
		p.victimSize = p.localSize
		p.local = nil // 主缓存置空
		p.localSize = 0
	}

	oldPools, allPools = allPools, nil
}

var (
	// 保护 allPools 的互斥锁
	allPoolsMu Mutex

	// allPools 是拥有非空主缓存（non-empty primary caches）的 pool 的集合
	// 保证并发安全的机制有两种：1) 通过 allPoolsMu 互斥锁和 pinning（固定）机制；2) 通过垃圾回收时的程序暂停 STW（Stop-The-World）
	allPools []*Pool

	// oldPools 是可能拥有非空 victim 缓存（non-empty victim caches）的 pool 的集合
	// 保证并发安全机制为 STW（Stop-The-World）
	oldPools []*Pool
)

func init() {
	// 将 poolCleanup 注册到 runtime，确保每次 GC 开始时自动被调用
	runtime_registerPoolCleanup(poolCleanup)
}

// 根据给定的索引 i，计算出指向 local 数组（[P]poolLocal）中第 i 个 poolLocal 元素的指针
func indexLocal(l unsafe.Pointer, i int) *poolLocal {
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*unsafe.Sizeof(poolLocal{}))
	return (*poolLocal)(lp)
}

// Implemented in runtime.
func runtime_registerPoolCleanup(cleanup func())
func runtime_procPin() int // 禁止当前 goroutine (G) 被调度器抢占，并将其固定到当前运行的逻辑处理器 (P) 上，同时返回该 P 的 ID
func runtime_procUnpin()   // 解除之前由 runtime_procPin 实施的固定，允许当前 G 再次被调度器抢占

// The below are implemented in internal/runtime/atomic and the
// compiler also knows to intrinsify the symbol we linkname into this
// package.

//go:linkname runtime_LoadAcquintptr internal/runtime/atomic.LoadAcquintptr
func runtime_LoadAcquintptr(ptr *uintptr) uintptr // 以 acquire 内存序加载一个 uintptr

//go:linkname runtime_StoreReluintptr internal/runtime/atomic.StoreReluintptr
func runtime_StoreReluintptr(ptr *uintptr, val uintptr) uintptr // 以 release 内存序存储一个 uintptr
