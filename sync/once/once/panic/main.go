package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		panic("Only once")
	}

	for i := 0; i < 5; i++ {
		func() {
			defer func() {
				r := recover()
				fmt.Println("recover", r)
			}()
			once.Do(onceBody)
		}()
	}
}
