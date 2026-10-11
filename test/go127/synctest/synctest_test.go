package synctest

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestTime(t *testing.T) {
	start := time.Now() // 起始时间为当前时间
	go func() {
		time.Sleep(1 * time.Second)
		t.Log(time.Since(start)) // 每次输出都是 "1s"
	}()
	time.Sleep(2 * time.Second)
	t.Log(time.Since(start)) // 每次输出都是 "2s"
}

func TestTimeWithSynctest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now() // 起始时间为 2000-01-01 UTC
		go func() {
			time.Sleep(1 * time.Second)
			t.Log(time.Since(start)) // 每次输出都是 "1s"
		}()
		time.Sleep(2 * time.Second)
		t.Log(time.Since(start)) // 每次输出都是 "2s"
	})
}
