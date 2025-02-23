package cron

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// JobWrapper 作业装饰器，可以为作业附加新功能
type JobWrapper func(Job) Job

// Chain 是存储 JobWrappers 的列表，它使用面向切面编程的横切行为（cross-cutting behaviors）来装饰提交的作业。
// 通过装饰链为作业附加各种新功能
type Chain struct {
	wrappers []JobWrapper
}

// NewChain 返回由给定 JobWrappers 组成的链
func NewChain(c ...JobWrapper) Chain {
	return Chain{c}
}

// Then 用链中的所有 JobWrappers 装饰给定的作业。
// 装饰器是顺序敏感的，设置不同的装饰器执行顺序，可能得到不同结果
//
// This:
//
//	NewChain(m1, m2, m3).Then(job)
//
// is equivalent to:
//
//	m1(m2(m3(job)))
func (c Chain) Then(j Job) Job {
	for i := range c.wrappers {
		j = c.wrappers[len(c.wrappers)-i-1](j)
	}
	return j
}

// Recover 捕获作业执行期间发生的 panic 并记录到日志中。
func Recover(logger Logger) JobWrapper {
	return func(j Job) Job {
		return FuncJob(func() {
			defer func() {
				if r := recover(); r != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					logger.Error(err, "panic", "stack", "...\n"+string(buf))
				}
			}()
			j.Run()
		})
	}
}

// DelayIfStillRunning 将作业串行化，延迟后续作业的执行，直到前一个作业完成。
// 如果作业的延迟超过一分钟，则会在记录 Info 级别的延迟日志。
func DelayIfStillRunning(logger Logger) JobWrapper {
	return func(j Job) Job {
		var mu sync.Mutex
		return FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.Info("delay", "duration", dur)
			}
			j.Run()
		})
	}
}

// SkipIfStillRunning 执行作业时，如果之前的调用仍在运行，则跳过此次作业的调用。
// 它会记录 Info 级别的跳过日志到给定的 logger。
func SkipIfStillRunning(logger Logger) JobWrapper {
	return func(j Job) Job {
		var ch = make(chan struct{}, 1)
		ch <- struct{}{}
		return FuncJob(func() {
			select {
			case v := <-ch:
				defer func() { ch <- v }()
				j.Run()
			default:
				logger.Info("skip")
			}
		})
	}
}
