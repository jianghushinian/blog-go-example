package synctest

import (
	"context"
	"testing"
	"testing/synctest"
)

func TestContextAfterFuncWithSynctest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// 创建一个可取消的 ctx
		ctx, cancel := context.WithCancel(t.Context())

		afterFuncCalled := false
		// context.AfterFunc 在 ctx 取消时会在自己的 goroutine 中调用 f
		context.AfterFunc(ctx, func() {
			afterFuncCalled = true
		})

		// 等待气泡内所有 goroutine 进入持久阻塞
		synctest.Wait()
		if afterFuncCalled { // 验证在 ctx 取消前 AfterFunc 没有被调用
			t.Fatalf("before context is canceled: AfterFunc called")
		}

		cancel() // 取消 ctx

		// 等待气泡内所有 goroutine 进入持久阻塞
		synctest.Wait()
		if !afterFuncCalled { // 验证在 ctx 取消后 AfterFunc 被调用
			t.Fatalf("after context is canceled: AfterFunc not called")
		}
	})
}
