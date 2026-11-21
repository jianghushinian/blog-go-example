package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		result int
		wg     sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			result++
		}()
	}
	wg.Wait()
	fmt.Println(result)
}
