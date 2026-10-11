//go:build concurrentsleep_demo

package synctest

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// 坑：synctest.Sleep = time.Sleep + synctest.Wait，而 Wait 在同一时刻只能有
// 一个调用者。多个 handler goroutine 同时 synctest.Sleep，会 panic
// "wait already in progress"，被 NewTestServer 捕获后报成 t.Errorf → FAIL。
//
// 正确做法：要在气泡内"等时间"的只有测试主流程（用 synctest.Wait/Sleep），
// handler 里用普通的 time.Sleep。
//
//	go test -tags concurrentsleep_demo -v .
func TestConcurrentSynctestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ts := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			synctest.Sleep(100 * time.Millisecond) // 错误示范：并发调用 Wait
			w.WriteHeader(http.StatusOK)
		}))
		client := ts.Client()

		var wg sync.WaitGroup
		for range 3 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := client.Get("http://example.com/")
				if err == nil {
					resp.Body.Close()
				}
			}()
		}
		wg.Wait()
	})
}
