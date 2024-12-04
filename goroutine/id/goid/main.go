package main

import (
	"fmt"
	"sync"

	"github.com/petermattis/goid"
)

func main() {
	fmt.Println("main", goid.Get())
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, goid.Get())
		}()
	}
	wg.Wait()
}

// $ go run goid/main.go
// main 1
// 9 43
// 4 38
// 5 39
// 6 40
// 7 41
// 8 42
// 1 35
// 0 34
// 2 36
// 3 37
