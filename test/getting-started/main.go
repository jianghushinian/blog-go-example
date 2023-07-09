package main

import (
	"fmt"

	"github.com/jianghushinian/blog-go-example/test/getting-started/abs"
	"github.com/jianghushinian/blog-go-example/test/getting-started/hello"
)

func main() {
	fmt.Println(abs.Abs(-1))
	fmt.Println(hello.Hello("World"))
}
