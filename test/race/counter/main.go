package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			counter++
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			counter++
		}
	}()
	wg.Wait()
	fmt.Println(counter)
}
