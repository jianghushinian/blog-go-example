package bar

import (
	_ "unsafe"
)

func add(a, b int) int {
	return a + b
}

// Push 模式（导出本地实现）

//go:linkname div github.com/jianghushinian/blog-go-example/directive/linkname/foo.Div
func div(a, b int) int {
	return a / b
}

// Handshake 模式（双方握手模式）

//go:linkname Hello
func Hello(name string) string {
	return "Hello " + name + "!"
}

//go:linkname x
var x = 12
