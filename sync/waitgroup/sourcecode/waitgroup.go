// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	// "internal/race"
	"sync/atomic"
	// "unsafe"
)

// WaitGroup 结构体
type WaitGroup struct {
	noCopy noCopy // 避免复制

	// 高 32 位是计数器（counter）的值，低 32 位是等待者（waiter）的数量
	state atomic.Uint64 // high 32 bits are counter, low 32 bits are waiter count.
	sema  uint32        // 信号量，用于 阻塞/唤醒 waiter
}

// Add 为计数器（counter）的值增加 delta（delta 可能为负数）
// 如果 counter 为 0，则唤醒所有被阻塞的 waiter
// 如果 counter 为负数，则 panic
func (wg *WaitGroup) Add(delta int) {
	// if race.Enabled {
	// 	if delta < 0 {
	// 		// Synchronize decrements with Wait.
	// 		race.ReleaseMerge(unsafe.Pointer(wg))
	// 	}
	// 	race.Disable()
	// 	defer race.Enable()
	// }
	state := wg.state.Add(uint64(delta) << 32) // delta 左移 32 位后与 state 相加，即为 counter 值加上 delta
	v := int32(state >> 32)                    // state 右移 32 位得到 counter 的值（这里拿到的是加上 delta 后的值）
	w := uint32(state)                         // state 转成 uint32 其实是直接拿低 32 位的值，得到 waiter 的值
	// if race.Enabled && delta > 0 && v == int32(delta) {
	// 	// The first increment must be synchronized with Wait.
	// 	// Need to model this as a read, because there can be
	// 	// several concurrent wg.counter transitions from 0.
	// 	race.Read(unsafe.Pointer(&wg.sema))
	// }
	if v < 0 { // 如果 counter 值为负数，直接 panic
		panic("sync: negative WaitGroup counter")
	}
	// 并发调用 Wait 和 Add 会触发 panic
	// - w != 0            表示有 waiter 存在，即调用过 Wait 方法，且还未被唤醒
	// - delta > 0         表示要增加 counter 的值（说明调用的肯定不是 Done 方法）
	// - v == int32(delta) 说明在调用 Add 方法之前，counter 值为 0
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// - v > 0  正常调用 Add 或 Done 方法
	// - w == 0 当前没有被阻塞的 waiter，即还未调用 Wait 方法
	if v > 0 || w == 0 { // 条件成立说明 counter 值加上 delta 操作成功，返回
		return
	}

	// 如果 counter 值为 0，并且还有被阻塞的 waiter，程序继续向下执行

	// 并发调用 Wait 和 Add 会触发 panic
	if wg.state.Load() != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// 目前 counter 值已经为 0，这里重置 waiter 数量为 0
	wg.state.Store(0)
	for ; w != 0; w-- { // 唤醒所有 waiter
		runtime_Semrelease(&wg.sema, false, 0)
	}
}

// Done 将计数器（counter）值减 1
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait 阻塞调用者当前的 goroutine（waiter），直到计数器（counter）值为 0
func (wg *WaitGroup) Wait() {
	// if race.Enabled {
	// 	race.Disable()
	// }
	for { // 使用无限循环保证 CAS 操作成功，因为并发调用是 CAS 操作可能失败需要重试
		state := wg.state.Load()
		v := int32(state >> 32) // 拿到 counter 值
		// w := uint32(state)   // 拿到 waiter 值
		if v == 0 { // 如果 counter 值已经为 0，直接返回
			// Counter is 0, no need to wait.
			// if race.Enabled {
			// 	race.Enable()
			// 	race.Acquire(unsafe.Pointer(wg))
			// }
			return
		}
		// 使用 CAS 操作增加 waiter 的数量
		if wg.state.CompareAndSwap(state, state+1) { // 将 waiter 数量 + 1，搭配外层的 for 循环，确保操作一定可以成功
			// if race.Enabled && w == 0 {
			// 	// Wait must be synchronized with the first Add.
			// 	// Need to model this is as a write to race with the read in Add.
			// 	// As a consequence, can do the write only for the first waiter,
			// 	// otherwise concurrent Waits will race with each other.
			// 	race.Write(unsafe.Pointer(&wg.sema))
			// }
			runtime_Semacquire(&wg.sema) // 阻塞当前 waiter 所在的 goroutine，等待被唤醒
			if wg.state.Load() != 0 {    // waiter 被唤醒后，此时 state 值理论上应该为 0（Add 方法中在唤醒 waiter 前会将其置 0：wg.state.Store(0)）
				// 如果 state 值不为 0，说明并发调用了 Wait 和 Add（重用 WaitGroup 时，当前批次 Wait 调用还未返回，就又在新批次调用了 Add），会触发 panic
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			// if race.Enabled {
			// 	race.Enable()
			// 	race.Acquire(unsafe.Pointer(wg))
			// }
			return // 如果 state 值为 0（计数器 counter 值为 0），说明 waiter 所等待的任务全部完成，成功返回
		}
	}
}
