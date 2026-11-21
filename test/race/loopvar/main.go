package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		i  int
		wg sync.WaitGroup
	)

	for i = 0; i < 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
		})
	}
	wg.Wait()
}
