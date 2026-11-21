package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				counter.Add(1)
			}
		})
	}
	wg.Wait()
	fmt.Println(counter.Load())
}
