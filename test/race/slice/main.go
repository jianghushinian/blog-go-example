package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		results []int
		wg      sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				results = append(results, j)
			}
		})
	}
	wg.Wait()
	fmt.Println(len(results))
}
