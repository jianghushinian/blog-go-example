// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package singleflight provides a duplicate function call suppression
// mechanism.
package singleflight // import "golang.org/x/sync/singleflight"

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// Sentinel error，用于标记在用户提供的 fn 函数中调用了 runtime.Goexit
// errGoexit indicates the runtime.Goexit was called in
// the user given function.
var errGoexit = errors.New("runtime.Goexit was called")

// A panicError is an arbitrary value recovered from a panic
// with the stack trace during the execution of given function.
type panicError struct {
	value interface{} // 记录 fn 函数的 panic 信息
	stack []byte      // 记录发生 panic 时的异常堆栈信息
}

// Error implements error interface.
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func (p *panicError) Unwrap() error {
	err, ok := p.value.(error)
	if !ok {
		return nil
	}

	return err
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// call is an in-flight or completed singleflight.Do call
type call struct {
	wg sync.WaitGroup // in-flight 并发控制

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val interface{} // 记录 fn 返回值
	err error       // 记录 fn 返回的 error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int             // 记录从缓存中获取 fn 返回值的次数
	chans []chan<- Result // 提供给 DoChan 方法用于传递 fn 的返回值
}

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// Result holds the results of Do, so they can be passed
// on a channel.
type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock() // 加锁，保证并发安全
	if g.m == nil {
		g.m = make(map[string]*call) // 延迟初始化 m
	}
	if c, ok := g.m[key]; ok { // 如果 key 已经在 map 中，即非第一个请求会进入到这个代码块
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call) // 当前 key 对应的第一个请求会创建一个 call 对象
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn) // 真正去执行 fn 的方法
	return c.val, c.err, c.dups > 0
}

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
//
// The returned channel will not be closed.
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1) // 构造一个 channel 用于传递 fn 的执行结果
	g.mu.Lock()                // 加锁，保证并发安全
	if g.m == nil {
		g.m = make(map[string]*call) // 延迟初始化 m
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}} // 创建一个 call 对象，并初始化 chans 字段
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn) // 开启新的 goroutine 来执行 fn

	return ch // 返回 channel 对象
}

// doCall handles the single call for a key.
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false // fn 是否正常返回
	recovered := false    // fn 是否产生 panic

	// 使用 double-defer 来区分 panic 或 runtime.Goexit
	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	defer func() {
		// 如果条件成立，则说明给定的函数 fn 内部调用了 runtime.Goexit
		// the given function invoked runtime.Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()        // 通知阻塞等待的其他请求可以获取 fn 执行结果了
		if g.m[key] == c { // fn 执行完成，从 m 中删除 key 记录
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			if len(c.chans) > 0 {
				go panic(e) // 为了防止等待 channel 的 goroutine 被永久阻塞，需要确保这个 panic 无法被 recover
				// 保持当前 goroutine 不退出
				select {} // Keep this goroutine around so that it will appear in the crash dump.
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// 当前 goroutine 正在执行 runtime.Goexit 退出流程，这里无需特殊处理
			// Already in the process of goexit, no need to call again
		} else {
			// 进入此代码块，说明 fn 正常返回
			// Normal return
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take a stack trace until we've determined
				// whether this is a panic or a runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				if r := recover(); r != nil { // 进入此代码块，说明 fn 触发了 panic
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (g *Group) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
