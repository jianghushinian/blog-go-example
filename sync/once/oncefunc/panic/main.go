package main

import (
	"fmt"
	"sync"
)

func main() {
	onceBody := sync.OnceFunc(func() {
		panic("Only once")
	})

	for i := 0; i < 5; i++ {
		func() {
			defer func() {
				r := recover()
				fmt.Println("recover", r)
			}()
			onceBody()
		}()
	}
}
