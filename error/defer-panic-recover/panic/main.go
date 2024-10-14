package main

import (
	"fmt"
	"time"
)

// func f() {
// 	defer fmt.Println("defer 1")
// 	fmt.Println(1)
// 	panic("woah")
// 	defer fmt.Println("defer 2")
// 	fmt.Println(2)
// }

func g() {
	fmt.Println("calling g")
	// 子 goroutine 中发生 panic，主 goroutine 也会退出
	go f(0)
	fmt.Println("called g")
}

func f(i int) {
	fmt.Println("panicking!")
	panic(fmt.Sprintf("i=%v", i))
	fmt.Println("printing in f", i) // 不会被执行
}

func main() {
	// f()

	g()
	time.Sleep(10 * time.Second)
}
