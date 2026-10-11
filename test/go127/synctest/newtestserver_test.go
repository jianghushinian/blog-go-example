package synctest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// NewTestServer 用内存网络，把 HTTP I/O 变成气泡内的 channel 通信，于是
// handler 里的 time.Sleep 走假时钟：10 个并发请求各耗时 100ms，整段测试
// 只在虚拟时间上前进 100ms，真实世界几乎不花时间。
//
// 注意 handler 里用的是 time.Sleep 而不是 synctest.Sleep，原因见
// concurrent_sleep_demo_test.go。
func TestNewTestServerWithSynctest(t *testing.T) {
	realStart := time.Now()

	synctest.Test(t, func(t *testing.T) {
		start := time.Now() // 气泡内虚拟时间

		ts := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			fmt.Fprint(w, "ok")
		}))
		client := ts.Client()

		var wg sync.WaitGroup
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := client.Get("http://example.com/")
				if err != nil {
					t.Errorf("request failed: %v", err)
					return
				}
				resp.Body.Close()
			}()
		}
		wg.Wait()

		if elapsed := time.Since(start); elapsed != 100*time.Millisecond {
			t.Fatalf("virtual elapsed = %v, want 100ms", elapsed)
		}
	})

	if real := time.Since(realStart); real >= 100*time.Millisecond {
		t.Fatalf("10 concurrent 100ms requests should cost ~no real time, took %v", real)
	}
}
