package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func main() {
	var (
		c  Counter
		wg sync.WaitGroup
	)

	for i := 0; i < 1000; i++ {
		wg.Go(func() {
			c.Inc()
		})
	}
	wg.Wait()
	fmt.Println(c.n)
}
