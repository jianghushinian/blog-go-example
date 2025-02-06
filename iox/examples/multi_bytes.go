package main

import (
	"fmt"

	"github.com/jianghushinian/blog-go-example/iox"
)

func main() {
	mb := iox.NewMultiBytes([][]byte{[]byte("Hello, World!\n")})
	_, _ = mb.Write([]byte("你好，世界！"))
	p := make([]byte, 32)
	_, _ = mb.Read(p)
	fmt.Println(string(p))
}

/*
$ go run examples/multi_bytes.go
Hello, World!
你好，世界！
*/
