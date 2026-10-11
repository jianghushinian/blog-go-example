package synctest

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// Wait 等待的是一条状态谓词：除调用者外，气泡内其它 goroutine 都处于持久阻塞状态或已退出。
func TestWaitSemantics(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var done atomic.Bool
		go func() {
			time.Sleep(1 * time.Second)
			done.Store(true)
		}()

		time.Sleep(1 * time.Second)
		t.Log(done.Load()) // 此处 done 可能仍然是 false
		synctest.Wait()
		t.Log(done.Load()) // 此处 done 必然是 true
	})
}
