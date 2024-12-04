package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func GoId() int {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func main() {
	fmt.Println("main", GoId())
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, GoId())
		}()
	}
	wg.Wait()
}

// $ go run main.go
// main 1
// 9 29
// 0 20
// 5 25
// 6 26
// 7 27
// 8 28
// 2 22
// 1 21
// 4 24
// 3 23
