package main

import (
	"fmt"

	"govanityurls.jianghushinian.cn/blog-go-example/doc/docgo/calculator"
)

func main() {
	fmt.Println(calculator.Add(5, 3))
	fmt.Println(calculator.Subtract(5, 3))
	fmt.Println(calculator.Multiply(5, 3))
	fmt.Println(calculator.Divide(6, 3))
}
