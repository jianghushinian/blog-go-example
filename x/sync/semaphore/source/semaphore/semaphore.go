// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package semaphore provides a weighted semaphore implementation.
package semaphore // import "golang.org/x/sync/semaphore"

import (
	"container/list"
	"context"
	"sync"
)

// 等待者结构体
type waiter struct {
	n     int64           // 请求资源数
	ready chan<- struct{} // 当获取到资源时被关闭，用于唤醒当前等待者
}

// NewWeighted 构造一个信号量对象
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

// Weighted 信号量结构体
type Weighted struct {
	size    int64      // 资源总数量
	cur     int64      // 当前已经使用的资源数
	mu      sync.Mutex // 互斥锁，保证对其他属性的操作并发安全
	waiters list.List  // 等待者队列，使用列表实现
}

// Acquire 请求 n 个资源
// 如果资源不足，则阻塞等待，直到有足够的资源数，或者 ctx 被取消
// 成功返回 nil，失败返回 ctx.Err() 并且不改变资源数
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	done := ctx.Done()

	s.mu.Lock() // 加锁保证并发安全

	// 如果在分配资源前 ctx 已经取消，则直接返回 ctx.Err()
	select {
	case <-done:
		s.mu.Unlock()
		return ctx.Err()
	default:
	}

	// 如果资源数足够，且不存在其他等待者，则请求资源成功，将 cur 加上 n，并返回
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n
		s.mu.Unlock()
		return nil
	}

	// 如果请求的资源数大于资源总数，不可能满足，则阻塞等待 ctx 取消，并返回 ctx.Err()
	if n > s.size {
		// Don't make other Acquire calls block on one that's doomed to fail.
		s.mu.Unlock()
		<-done
		return ctx.Err()
	}

	// 资源不够或者存在其他等待者，则继续执行

	// 加入等待队列
	ready := make(chan struct{})    // 创建一个 channel 作为一个属性记录到等待者对象 waiter 中，用于后续通知其唤醒
	w := waiter{n: n, ready: ready} // 构造一个等待者对象 waiter
	elem := s.waiters.PushBack(w)   // 将 waiter 追加到等待者队列
	s.mu.Unlock()

	// 使用 select 实现阻塞等待
	select {
	case <-done: // 检查 ctx 是否被取消
		s.mu.Lock()
		select {
		case <-ready: // 检查当前 waiter 是否被唤醒
			// 进入这里，说明是 ctx 被取消后 waiter 被唤醒
			s.cur -= n        // 那么就当作 waiter 没有被唤醒，将请求的资源数还回去
			s.notifyWaiters() // 通知等待队列，检查队列中下一个 waiter 资源数是否满足
		default: // 将当前 waiter 从等待者队列中移除
			isFront := s.waiters.Front() == elem // 当前 waiter 是否为第一个等待者
			s.waiters.Remove(elem)               // 从队列中移除
			// 如果当前 waiter 是队列中第一个等待者，并且还有剩余的资源数
			if isFront && s.size > s.cur {
				s.notifyWaiters() // 通知等待队列，检查队列中下一个 waiter 资源数是否满足
			}
		}
		s.mu.Unlock()
		return ctx.Err()

	case <-ready: // 被唤醒
		select {
		case <-done: // 再次检查 ctx 是否被取消
			// 进入这里，说明 waiter 被唤醒后 ctx 却被取消了，当作未被唤醒来处理
			s.Release(n) // 释放资源
			return ctx.Err()
		default:
		}
		return nil // 成功返回 nil
	}
}

// TryAcquire 尝试请求 n 个资源
// 不阻塞，成功返回 true，失败返回 false 并且不改变资源数
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock() // 加锁保证并发安全
	// 剩余资源数足够，且不存在其他等待者，则请求资源成功
	success := s.size-s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n // 记录当前已经使用的资源数
	}
	s.mu.Unlock()
	return success
}

// Release 释放 n 个资源
func (s *Weighted) Release(n int64) {
	s.mu.Lock() // 加锁保证并发安全
	s.cur -= n  // 释放资源
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: released more than held")
	}
	s.notifyWaiters() // 通知等待队列，检查队列中下一个 waiter 资源数是否满足
	s.mu.Unlock()
}

// 检查队列中下一个 waiter 资源数是否满足
func (s *Weighted) notifyWaiters() {
	// 循环检查下一个 waiter 请求的资源数是否满足，满足则出队，不满足则终止循环
	for {
		next := s.waiters.Front() // 获取队首元素
		if next == nil {
			break // 没有 waiter 了，队列为空终止循环
		}

		w := next.Value.(waiter)
		if s.size-s.cur < w.n { // 当前 waiter 资源数不满足，退出循环
			// 不继续查找队列中后续 waiter 请求资源是否满足，避免产生饥饿
			break
		}

		// 资源数满足，唤醒 waiter
		s.cur += w.n           // 记录使用的资源数
		s.waiters.Remove(next) // 从队列中移除 waiter
		close(w.ready)         // 利用关闭 channel 的操作，来唤醒 waiter
	}
}
