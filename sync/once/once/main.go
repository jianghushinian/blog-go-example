package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

// func main() {
// 	var once sync.Once
// 	onceBody := func() {
// 		fmt.Println("Only once")
// 	}
// 	for i := 0; i < 10; i++ {
// 		once.Do(onceBody)
// 	}
// }

// func main() {
// 	var once sync.Once
// 	var i = 10
// 	onceBody := func() {
// 		i *= 2
// 	}
// 	done := make(chan bool)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			once.Do(onceBody)
// 			done <- true
// 		}()
// 	}
// 	for i := 0; i < 10; i++ {
// 		<-done
// 	}
// 	fmt.Println("i", i)
// }
