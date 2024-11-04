// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errgroup 提供了同步、错误传播和上下文取消功能，用于一组协程处理共同任务的子任务
//
// [errgroup.Group] 与 [sync.WaitGroup] 相关，但增加了处理任务返回错误的能力
package errgroup

import (
	"context"
	"fmt"
	"sync"
)

// 定义一个空结构体类型 token，会作为信号进行传递，用于控制并发数
type token struct{}

// Group 是一组协程的集合，这些协程处理同一整体任务的子任务
//
// 零值 Group 是有效的，对活动协程的数量没有限制，并且不会在出错时取消
type Group struct {
	cancel func(error) // 取消函数，就是 context.CancelCauseFunc 类型

	wg sync.WaitGroup // 内部使用了 sync.WaitGroup

	sem chan token // 信号 channel，可以控制协程并发数量

	errOnce sync.Once // 确保错误仅处理一次
	err     error     // 记录子协程集中返回的第一个错误
}

// 当一个协程完成时，调用此方法
func (g *Group) done() {
	// 如果设置了最大并发数，则 sem 不为 nil，从 channel 中消费一个 token，表示一个协程已完成
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done() // 转发给 sync.WaitGroup.Done()，将活动协程数减一
}

// WithContext 返回一个新的 Group 和一个从 ctx 派生的关联 Context
//
// 派生的 Context 会在传递给 Go 的函数首次返回非 nil 错误或 Wait 首次返回时被取消，以先发生者为准。
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := withCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait 会阻塞，直到来自 Go 方法的所有函数调用返回，然后返回它们中的第一个非 nil 错误（如果有的话）
func (g *Group) Wait() error {
	g.wg.Wait()          // 转发给 sync.WaitGroup.Wait()，等待所有协程执行完成
	if g.cancel != nil { // 如果 cancel 不为 nil，则调用取消函数，并设置 cause
		g.cancel(g.err)
	}
	return g.err // 返回错误
}

// Go 会在新的协程中调用给定的函数
// 它会阻塞，直到可以在不超过配置的活跃协程数量限制的情况下添加新的协程
//
// 首次返回非 nil 错误的调用会取消该 Group 的上下文（context），如果该 context 是通过调用 WithContext 创建的，该错误将由 Wait 返回
func (g *Group) Go(f func() error) {
	if g.sem != nil { // 这个是限制并发数的信号通道
		g.sem <- token{} // 如果超过了配置的活跃协程数量限制，向 channel 发送 token 会阻塞
	}

	g.wg.Add(1) // 转发给 sync.WaitGroup.Add(1)，将活动协程数加一
	go func() {
		defer g.done() // 当一个协程完成时，调用此方法，内部会将调用转发给 sync.WaitGroup.Done()

		if err := f(); err != nil { // f() 就是我们要执行的任务
			g.errOnce.Do(func() { // 仅执行一次，即只处理一次错误，所以会记录第一个非 nil 的错误，与协程启动顺序无关
				g.err = err          // 记录错误
				if g.cancel != nil { // 如果 cancel 不为 nil，则调用取消函数，并设置 cause
					g.cancel(g.err)
				}
			})
		}
	}()
}

// TryGo 仅在 Group 中活动的协程数量低于限额时，才在新的协程中调用给定的函数
//
// 返回值标识协程是否启动
func (g *Group) TryGo(f func() error) bool {
	if g.sem != nil { // 如果设置了最大并发数
		select {
		case g.sem <- token{}: // 可以向 channel 写入 token，说明没有达到限额，可以启动协程
			// Note: this allows barging iff channels in general allow barging.
		default: // 如果超过了配置的活跃协程数量限制，会走到这个 case
			return false
		}
	}

	// 接下来的代码与 Go 中的逻辑相同
	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
	return true
}

// SetLimit 限制该 Group 中活动的协程数量最多为 n，负值表示没有限制
//
// 任何后续对 Go 方法的调用都将阻塞，直到可以在不超过限额的情况下添加活动协程
//
// 在 Group 中存在任何活动的协程时，限制不得修改
func (g *Group) SetLimit(n int) { // 传进来的 n 就是 channel 长度，以此来限制协程的并发数量
	if n < 0 { // 这里检查如果小于 0 则不限制协程并发数量。此外，也不要将其设置为 0，会产生死锁
		g.sem = nil
		return
	}
	if len(g.sem) != 0 { // 如果存在活动的协程，调用此方法将产生 panic
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	g.sem = make(chan token, n)
}
