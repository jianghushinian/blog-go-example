package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		mu sync.Mutex
		m  = make(map[int]int)
		wg sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				mu.Lock()
				m[j] = j
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	fmt.Println(len(m))
}
