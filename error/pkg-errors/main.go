package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func a() error {
	// 初始错误
	return errors.New("a error")
}

func b() error {
	err := a()
	if err != nil {
		// 包装新的错误并返回
		newErr := errors.WithMessage(err, "b error")
		// newErr := errors.Wrap(err, "b error")

		// 可以从包装后的错误中还原出初始错误
		fmt.Printf("newErr cause == err: %t\n", errors.Cause(newErr) == err)

		return newErr
	}
	return nil
}

func main() {
	err := b()
	if err != nil {
		// %v 打印错误信息
		fmt.Printf("%v\n", err)

		fmt.Println("============================================")

		// %+v 打印错误信息和错误堆栈
		fmt.Printf("%+v\n", err)

		fmt.Println("============================================")

		// 打印错误根因
		fmt.Printf("%v\n", errors.Cause(err))
		return
	}
	fmt.Println("success")
}

/*
$ go run main.go
newErr cause == err: true
b error: a error
============================================
a error
main.a
        /go/blog-go-example/error/pkg-errors/main.go:11
main.b
        /go/blog-go-example/error/pkg-errors/main.go:15
main.main
        /go/blog-go-example/error/pkg-errors/main.go:30
runtime.main
        /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/proc.go:271
runtime.goexit
        /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/asm_arm64.s:1222
b error
============================================
a error
*/
