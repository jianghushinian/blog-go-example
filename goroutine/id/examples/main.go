package main

import (
	"fmt"
	"runtime"
)

func main() {
	go func() {
		fmt.Println("goroutine running")
	}()

	buf := make([]byte, 1024)
	// n := runtime.Stack(buf, true)
	n := runtime.Stack(buf, false)
	fmt.Printf("written: %d\n", n)
	fmt.Printf("stack:\n%s\n", buf)
}
