package sync

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done atomic.Uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if o.done.Load() == 0 {
		o.m.Lock()
		defer o.m.Unlock()
		if o.done.Load() == 0 {
			defer o.done.Store(1)
			f()
		}
	}
}
