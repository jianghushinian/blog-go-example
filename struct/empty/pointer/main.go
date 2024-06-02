package main

import "fmt"

// NOTE: 实验查看空结构体内存地址是否相同

func main() {
	var (
		a struct{}
		b struct{}
		c struct{}
		d struct{}
	)

	println("&a:", &a)
	println("&b:", &b)
	println("&c:", &c)
	println("&d:", &d)

	println("&a == &b:", &a == &b)
	x := &a
	y := &b
	println("x == y:", x == y)

	fmt.Printf("&c(%p) == &d(%p): %t\n", &c, &d, &c == &d)
}

// $ go run -gcflags='-m -N -l' main.go
// # command-line-arguments
// ./main.go:11:3: moved to heap: c
// ./main.go:12:3: moved to heap: d
// ./main.go:23:12: ... argument does not escape
// ./main.go:23:50: &c == &d escapes to heap
// &a: 0x1400010ae84
// &b: 0x1400010ae84
// &c: 0x104ec74a0
// &d: 0x104ec74a0
// &a == &b: false
// x == y: true
// &c(0x104ec74a0) == &d(0x104ec74a0): true
