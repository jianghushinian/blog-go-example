package main

import "fmt"

func main() {
	var x int
	ch := make(chan struct{})

	go func() {
		x = 1
		ch <- struct{}{}
	}()

	<-ch
	fmt.Println(x)
}
