// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync/atomic"
)

// Map 支持并发操作的 map
type Map struct {
	// 互斥锁，用于保护对 dirty 字段的并发访问
	// 需要修改 dirty 或在 read 和 dirty 之间同步数据时，必须持有该锁
	// 其他一些操作（如从 dirty 提升到 read）需要借助这个锁来确保线程安全
	mu Mutex

	// read 包含了可以安全进行并发访问的部分 map 内容（无论是否持有 mu）。
	//
	// read 字段本身始终可以安全读取，但只有在持有 mu（互斥锁）时才能进行写入。
	//
	// 存储在 read 中的条目可以在不持有 mu 的情况下被并发更新，
	// 但如果要更新一个先前被删除（expunged）的条目，
	// 则必须先将该条目复制到 dirty map，并在持有 mu 的情况下取消删除（unexpunged）。
	read atomic.Pointer[readOnly]

	// dirty 包含了需要持有 mu（互斥锁）才能访问的部分 map 内容。
	// 为了确保 dirty map 可以快速提升为 read map，
	// 它还包含了 read map 中所有未被删除（未标记为 expunged）的条目。
	//
	// 被删除（expunged）的条目不会存储在 dirty map 中。
	// 在 clean map 中被删除的条目，必须先取消删除（unexpunged）并添加到 dirty map，
	// 然后才能存入新的值。
	//
	// 如果 dirty map 为 nil，则下一次写入 map 时，会通过对 clean map 进行浅拷贝（省略过时的条目）来初始化 dirty map。
	dirty map[any]*entry

	// misses 统计自上次更新 read map 以来，
	// 由于无法直接在 read map 中确定键是否存在，而需要加锁 (mu) 进行查找的次数。
	//
	// 一旦累计的 misses 次数足以实现复制 dirty map 的开销，
	// dirty map 就会被提升为 read map（处于未修改状态），
	// 并且下一次写入 map 时，会创建一个新的 dirty map 副本。
	misses int
}

// readOnly 是一个不可变结构体，它被原子地存储在 Map.read 字段中。
type readOnly struct {
	m       map[any]*entry // 存储键值对的只读 map。
	amended bool           // 如果 dirty map 中包含 read map (m) 中不存在的键，则为 true。
}

// expunged 是一个任意的指针地址，用于标记从 dirty map 中被删除的条目。
var expunged = new(any)

// entry 代表 map 中与特定键对应的存储槽。
// 其实就是 map 中键对应的值。
type entry struct {
	// p 指向存储在该条目中的 interface{} 值（任意类型的值）。
	//
	// 如果 p == nil，表示该条目已被删除，并且要么 m.dirty == nil，要么 m.dirty[key] 指向该 entry。
	//
	// 如果 p == expunged，表示该条目已被彻底删除，且 m.dirty 不是 nil，同时该条目不会出现在 m.dirty 中。
	//
	// 否则，该条目是有效的，并记录在 m.read.m[key] 中，如果 m.dirty 不是 nil，则也记录在 m.dirty[key] 中。
	//
	// 可以通过原子替换为 nil 来删除一个条目：当下次创建 m.dirty 时，它会将 nil 原子地替换为 expunged，
	// 并且不会在 m.dirty[key] 中设置该条目。
	//
	// 只要 p != expunged，就可以通过原子替换来更新条目关联的值。
	// 如果 p == expunged，则必须先在 m.dirty[key] 中存入该 entry，
	// 以便 dirty map 的查找操作能够找到该条目，然后才能更新其值。
	p atomic.Pointer[any]
}

// 给定 map 的 value 创建一个条目 entry
func newEntry(i any) *entry {
	e := &entry{}
	e.p.Store(&i)
	return e
}

// 加载只读 map
func (m *Map) loadReadOnly() readOnly {
	if p := m.read.Load(); p != nil {
		return *p
	}
	return readOnly{}
}

// Load 返回 map 中存储的指定键（key）对应的值，如果 key 不存在，则返回 nil。
// 返回值 ok 表示是否在 map 中找到了该 key 对应的值。
func (m *Map) Load(key any) (value any, ok bool) {
	read := m.loadReadOnly() // 加载只读 map
	e, ok := read.m[key]     // 尝试从 read map 中查找 key

	// 如果 read map 中找不到 key，并且 dirty map 可能包含新键（amended 为 true），则进入慢路径
	if !ok && read.amended {
		m.mu.Lock() // 操作 dirty 需要加锁
		// 在获取锁的过程中，dirty 可能已经被提升为 read，因此需要重新加载 read map 进行 double-checking
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			// 如果在 read map 仍然找不到，则尝试从 dirty map 查找
			e, ok = m.dirty[key] // 从 dirty map 中查找 key
			// 无论 key 是否存在，记录一次 miss：
			// 在 dirty map 提升为 read map 之前，这个 key 的查询都会走慢路径
			m.missLocked() // miss 计数值 + 1
		}
		m.mu.Unlock()
	}

	// 如果最终 key 仍然不存在，则返回 nil
	if !ok {
		return nil, false
	}
	return e.load() // 返回找到的值
}

// 加载 entry 中存储的值。
func (e *entry) load() (value any, ok bool) {
	p := e.p.Load() // 原子加载指针 p
	// 已经被删除
	if p == nil || p == expunged {
		return nil, false
	}
	return *p, true
}

// Store 为指定的键（key）设置值。
func (m *Map) Store(key, value any) {
	// 调用 Swap 方法，将 key 关联的值替换为新值，并忽略旧值
	_, _ = m.Swap(key, value)
}

// Clear 删除所有条目，使 Map 变为空。
func (m *Map) Clear() {
	read := m.loadReadOnly() // 加载只读 map
	// 如果 read map 为空且 dirty map 也未修改（amended 为 false），则无需执行清理
	if len(read.m) == 0 && !read.amended {
		// 如果 map 已经为空，则避免分配新的 readOnly。
		return
	}

	m.mu.Lock() // 加锁
	defer m.mu.Unlock()

	// double-checking
	// 重新加载 read map，避免在获取锁的过程中发生变化
	read = m.loadReadOnly()
	if len(read.m) > 0 || read.amended {
		// 清空 read map，存储一个新的空 readOnly 结构
		m.read.Store(&readOnly{})
	}

	// 清空 dirty map
	clear(m.dirty)
	// 复位 miss 计数，防止下一次操作时刚清空的 dirty map 立即被提升到 read map
	m.misses = 0
}

// tryCompareAndSwap 比较 entry 中的值是否等于给定的 old 值，
// 如果相等且条目未被删除（未标记为 expunged），则将其替换为 new 值。
//
// 如果条目已被删除（expunged），tryCompareAndSwap 返回 false，且不对 entry 做任何修改。
func (e *entry) tryCompareAndSwap(old, new any) bool {
	// 加载当前 entry 存储的值
	p := e.p.Load()
	// 如果值为 nil（条目已删除）或 expunged（条目已被彻底删除）或当前值不等于 old，则返回 false
	if p == nil || p == expunged || *p != old {
		return false
	}

	// 复制 new 以优化逃逸分析，避免在比较失败时不必要的堆分配
	nc := new
	// 循环执行 CAS 操作，直到成功返回
	for {
		// 尝试使用原子操作 CompareAndSwap 替换值
		if e.p.CompareAndSwap(p, &nc) {
			return true
		}

		// 每次循环都要执行一次检查
		// 重新加载最新的值，确保并发情况下仍满足条件
		p = e.p.Load()
		// 如果值发生变化，或者条目已删除，则返回 false
		if p == nil || p == expunged || *p != old {
			return false
		}
	}
}

// unexpungeLocked 确保 entry 未被标记为 expunged（彻底删除）。
//
// 如果 entry 之前被标记为 expunged，则在 m.mu 解锁之前，该 entry 必须被添加回 dirty map。
func (e *entry) unexpungeLocked() (wasExpunged bool) {
	// 使用 CompareAndSwap 将 expunged 替换为 nil，如果替换成功，说明之前是 expunged
	return e.p.CompareAndSwap(expunged, nil)
}

// swapLocked 无条件地将一个新值存入 entry。
//
// 调用此方法前必须确保 entry 未被 expunged。
func (e *entry) swapLocked(i *any) *any {
	// 使用原子 Swap 操作替换存储的值，并返回旧值
	return e.p.Swap(i)
}

// LoadOrStore 如果 key 存在，则返回其当前值；否则，存储并返回给定的 value。
// loaded 返回值表示是否是从 map 中加载的值（true 表示 key 已存在，false 表示存储了新值）。
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool) {
	// 优先检查只读 map，避免加锁，提高并发性能
	read := m.loadReadOnly() // 加载只读 map
	if e, ok := read.m[key]; ok {
		// 尝试加载或存储值
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	// 加锁处理
	m.mu.Lock()
	defer m.mu.Unlock()

	// 进入慢路径

	// 重新加载 read map 以防数据变化
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok { // key 在 read map 中
		// 解除 expunged 状态并放入 dirty map
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		// 尝试加载或存储
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok { // key 仅存在于 dirty map 中
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked() // 记录一次 miss
	} else { // 该 key 完全不存在
		if !read.amended {
			// 这是 dirty map 第一次存入新 key，需要确保它已分配，并标记 read map 为不完整
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		// 在 dirty map 中存储新值
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}

	return actual, loaded
}

// tryLoadOrStore 在 entry 未被 expunged（彻底删除）时，原子地加载或存储一个值。
//
// 如果 entry 已被 expunged，tryLoadOrStore 不会修改 entry，并返回 ok==false。
func (e *entry) tryLoadOrStore(i any) (actual any, loaded, ok bool) {
	p := e.p.Load()
	if p == expunged {
		// entry 被标记为 expunged，不可修改
		return nil, false, false
	}
	if p != nil {
		// entry 已经存在，loaded 为 true，返回当前值
		return *p, true, true
	}

	// 复制 i，以提高逃逸分析的优化效果：
	// 如果 hit 了 "load" 分支或者 entry 被 expunged，就不需要额外的堆分配。
	ic := i
	// 循环执行 CAS 操作，直到成功返回
	for {
		// 如果当前值为 nil，则尝试存储新值
		if e.p.CompareAndSwap(nil, &ic) {
			return i, false, true // 存储成功，loaded 为 false
		}

		// 每次循环都要执行一次检查
		p = e.p.Load()
		if p == expunged {
			// entry 被 expunged，返回失败
			return nil, false, false
		}
		if p != nil {
			// entry 已经存在，loaded 为 true，返回当前值
			return *p, true, true
		}
	}
}

// LoadAndDelete 删除 key 对应的值，并返回删除前的值（如果存在）。
// loaded 返回该 key 是否存在。
func (m *Map) LoadAndDelete(key any) (value any, loaded bool) {
	// 先从 read map 读取，避免加锁
	read := m.loadReadOnly() // 加载只读 map
	e, ok := read.m[key]
	if !ok && read.amended {
		// 若 read map 没有 key，且 dirty map 可能包含该 key，则加锁检查 dirty map
		m.mu.Lock()
		read = m.loadReadOnly() // double-checking
		e, ok = read.m[key]
		if !ok && read.amended {
			// 在 dirty map 查找 key
			e, ok = m.dirty[key]
			delete(m.dirty, key) // 删除 dirty map 中的 key
			// 记录一次 miss，表示该 key 走了慢路径，直到 dirty map 提升为 read map
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		// 若 key 存在，调用 entry.delete() 删除值并返回
		return e.delete()
	}
	return nil, false // read map 和 dirty map 都无此 key
}

// Delete 删除 key 对应的值。
func (m *Map) Delete(key any) {
	m.LoadAndDelete(key) // 直接调用 LoadAndDelete 进行删除，忽略其返回值
}

// 删除 entry 存储的值，并返回删除前的值和是否成功删除的标志。
// 这里的删除只会将 entry 中存储的值置为 nil，并不会彻底删除 entry 对象
func (e *entry) delete() (value any, ok bool) {
	// 循环执行 CAS 操作，直到成功返回
	for {
		p := e.p.Load()
		// 如果 entry 为空（已删除）或者已经被标记为 expunged（彻底删除），返回 false
		if p == nil || p == expunged {
			return nil, false
		}
		// 尝试原子替换，将存储的值置为 nil，删除成功则返回原值
		if e.p.CompareAndSwap(p, nil) {
			return *p, true
		}
	}
}

// trySwap 交换 entry 中的值，如果 entry 没有被标记为 expunged（彻底删除）。
//
// 如果 entry 已被 expunged，则返回 false 并保持 entry 不变。
func (e *entry) trySwap(i *any) (*any, bool) {
	// 循环执行 CAS 操作，直到成功返回
	for {
		p := e.p.Load()
		// 如果 entry 已被 expunged，则无法替换，返回 nil 和 false
		if p == expunged {
			return nil, false
		}
		// 尝试原子替换，如果成功，则返回原值和 true
		if e.p.CompareAndSwap(p, i) {
			return p, true
		}
		// 如果 CompareAndSwap 失败，则说明 p 被其他线程修改了，需要重新尝试
	}
}

// Swap 交换 key 对应的值，并返回之前的值（如果存在）。
// 返回值 loaded 表示 key 是否存在。
func (m *Map) Swap(key, value any) (previous any, loaded bool) {
	// 先尝试从 read map 获取 key 对应的 entry
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok { // 如果 key 存在
		// 尝试交换值
		if v, ok := e.trySwap(&value); ok { // 如果交换成功
			if v == nil { // 说明已被删除，但没有彻底删除（expunged）
				return nil, false // 新增键值对成功
			}
			return *v, true // 交换成功
		}
	}

	// 未找到或需要修改 dirty map，需要加锁处理
	m.mu.Lock()
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok { // 如果 key 在 read map 中
		if e.unexpungeLocked() {
			// 如果 entry 之前被 expunged，意味着 dirty map 不为 nil 且该 entry 不在 dirty map 中
			m.dirty[key] = e
		}
		// 进行值交换
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else if e, ok := m.dirty[key]; ok { // 如果 key 在 dirty map 中
		// 进行值交换
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else { // key 即不在 read map 中，也不在 dirty map 中
		if !read.amended {
			// 添加第一个新 key 到 dirty map，需要先标记 read map 为不完整
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		// 在 dirty map 中创建新 entry
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
	return previous, loaded
}

// CompareAndSwap 仅在 key 对应的值等于 old 时，将其替换为 new。
// old 必须是可比较的类型。
func (m *Map) CompareAndSwap(key, old, new any) (swapped bool) {
	// 先尝试从 read map 查找 key
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok { // 如果 key 在 read map 中
		// 直接在 entry 上尝试 CompareAndSwap
		return e.tryCompareAndSwap(old, new)
	} else if !read.amended {
		// 如果 key 不存在，并且 read map 没有标记为 amended，说明 dirty map 也不会有 key
		return false
	}

	// 进入 dirty map 处理，需要加锁
	m.mu.Lock()
	defer m.mu.Unlock()
	read = m.loadReadOnly()

	if e, ok := read.m[key]; ok { // 如果 key 在 read map 中
		swapped = e.tryCompareAndSwap(old, new)
	} else if e, ok := m.dirty[key]; ok { // 如果 key 在 dirty map 中
		swapped = e.tryCompareAndSwap(old, new)
		// 这里虽然 CompareAndSwap 不会修改 map 的 key 集合，
		// 但因为涉及锁的操作，需要增加一次 miss 计数，
		// 这样 dirty map 之后会被提升为 read，提高访问效率。
		m.missLocked()
	}
	return swapped
}

// CompareAndDelete 仅在 key 存在且值等于 old 时删除该 entry。
// old 必须是可比较的类型。
//
// 如果 key 在 map 中不存在，则返回 false（即使 old 为 nil）。
func (m *Map) CompareAndDelete(key, old any) (deleted bool) {
	// 先尝试从 read map 查找 key
	read := m.loadReadOnly()
	e, ok := read.m[key]

	// 如果 key 不在 read map 中，但 read.amended == true，说明可能在 dirty map
	if !ok && read.amended {
		m.mu.Lock() // 进入 dirty map 处理，需要加锁
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 不直接从 m.dirty 删除 key，而是仅执行“compare”逻辑。
			// 当 dirty map 未来被提升为 read 时，该 entry 会被标记 expunged。
			//
			// 记录一次 miss，这样该 key 未来会走慢路径，直到 dirty map 变为 read。
			m.missLocked()
		}
		m.mu.Unlock()
	}

	// 如果 key 存在于 read map 或 dirty map
	for ok {
		p := e.p.Load()
		// 如果值已被删除（nil）、expunged（彻底删除）或者值不等于 old，则返回 false
		if p == nil || p == expunged || *p != old {
			return false
		}
		// CAS 交换值为 nil，成功删除
		if e.p.CompareAndSwap(p, nil) {
			return true
		}
	}
	return false
}

// Range 遍历 Map 中的所有键值对，并依次调用 f(key, value)。
// 如果 f 返回 false，则停止遍历。
//
// 需要注意的是，Range 并不保证访问的是某个固定时刻的快照：
// - 每个 key 只会被访问一次。
// - 但如果遍历过程中有并发修改 key 的值或删除 key，则可能访问不同时间点的值。
// - f 函数内部可以对 Map 进行任何操作（包括插入、删除等）。
//
// 如果 f 只执行常数次操作，整个 Range 仍可能是 O(N) 复杂度。
func (m *Map) Range(f func(key, value any) bool) {
	// 读取当前的 read map
	read := m.loadReadOnly()

	// 如果 read map 是完整的（未修改），直接遍历
	if read.amended {
		// 若 read.amended == true，说明 dirty map 存在新增 key，需要提升为 read
		m.mu.Lock()
		read = m.loadReadOnly()
		if read.amended {
			// 提升 dirty map 为 read
			read = readOnly{m: m.dirty}
			copyRead := read
			m.read.Store(&copyRead)
			// 清空 dirty map 并重置 miss 计数
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	// 遍历 read map
	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		// 调用用户提供的函数 f
		if !f(k, v) {
			break // f 返回 false 时停止遍历
		}
	}
}

// NOTE: 调用 xxxLocked 方法之前需要持有锁

// 增加 misses 计数：
// misses 记录了访问 read map 但未找到 key 的次数。
func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) { // misses 计数值小于 dirty 数据总量，直接返回
		return
	}
	// misses 计数值大于等于 dirty 数据总量，则将 dirty 提升为 read，直接替换对象
	m.read.Store(&readOnly{m: m.dirty})
	m.dirty = nil // 标记 dirty 为 nil
	m.misses = 0  // misses 计数清零
}

// 当 dirty map 为空时，将 read map 拷贝到 dirty map，为未来的写操作做准备
// 最坏时间复杂度 O(N)，会导致性能劣化
func (m *Map) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read := m.loadReadOnly()
	m.dirty = make(map[any]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() { // 将未标记为彻底删除的 entry 记录到 dirty map 中
			m.dirty[k] = e
		}
	}
}

// 尝试将 entry 从 nil 变为 expunged，即标记为 "彻底删除"。
func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := e.p.Load()
	for p == nil {
		if e.p.CompareAndSwap(nil, expunged) {
			return true
		}
		p = e.p.Load()
	}
	return p == expunged
}
