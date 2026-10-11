package synctest

import (
	"testing"
	"testing/synctest"
	"time"
)

// Go 1.27: synctest.Sleep combines time.Sleep and synctest.Wait.
func TestSleep(t *testing.T) {
	realStart := time.Now()

	synctest.Test(t, func(t *testing.T) {
		start := time.Now() // virtual time inside the bubble

		woke := make(chan time.Duration, 1)
		go func() {
			time.Sleep(time.Second)
			woke <- time.Since(start)
		}()

		synctest.Sleep(2 * time.Second)

		select {
		case d := <-woke:
			if d != time.Second {
				t.Fatalf("worker virtual sleep = %v, want 1s", d)
			}
		default:
			t.Fatal("worker goroutine not settled after synctest.Sleep")
		}

		if elapsed := time.Since(start); elapsed != 2*time.Second {
			t.Fatalf("virtual elapsed = %v, want 2s", elapsed)
		}
	})

	if real := time.Since(realStart); real >= 100*time.Millisecond {
		t.Fatalf("2s of virtual time should cost ~no real time, took %v", real)
	}
}
