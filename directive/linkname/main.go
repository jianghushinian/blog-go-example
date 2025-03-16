package main

import (
	"fmt"

	"github.com/jianghushinian/blog-go-example/directive/linkname/foo"

	_ "unsafe"
)

//go:linkname TooLarge 	fmt.tooLarge
func TooLarge(x int) bool

func main() {
	fmt.Println("foo.Add(1, 2):", foo.Add(1, 2))
	fmt.Println("foo.Div(2, 1):", foo.Div(2, 1))
	fmt.Println(`foo.Hello("jianghushinian"):`, foo.Hello("jianghushinian"))
	fmt.Println("foo.Now():", foo.Now())
	foo.FooPrintln("Calling FooPrintln")
	fmt.Println("foo.X:", foo.X)
	fmt.Println("foo.TooLarge(1e6 + 1):", foo.TooLarge(1e6+1))
	// ref: https://github.com/golang/go/issues/67401
	// 执行时需要使用：go run -ldflags=-checklinkname=0 main.go
	// 否则会报错：link: main: invalid reference to fmt.tooLarge
	fmt.Println("TooLarge(1e6 + 1):", TooLarge(1e6+1))
}
