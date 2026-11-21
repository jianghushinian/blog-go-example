package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var (
		result atomic.Int64
		wg     sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Go(func() {
			result.Add(1)
		})
	}
	wg.Wait()
	fmt.Println(result.Load())
}
