package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		mu      sync.Mutex
		results []int
		wg      sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				mu.Lock()
				results = append(results, j)
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	fmt.Println(len(results))
}
