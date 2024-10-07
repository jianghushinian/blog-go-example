package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// NOTE: 记录错误调用链
//
// func Foo() error {
// 	return errors.New("foo error")
// }
//
// func Bar() error {
// 	err := Foo()
// 	if err != nil {
// 		return errors.Wrap(err, "bar")
// 	}
// 	return nil
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		fmt.Printf("err: %s\n", err)
// 	}
// }

// NOTE: 记录错误堆栈

// func Foo() error {
// 	return errors.New("foo error")
// }
//
// // func Bar() error {
// // 	err := Foo()
// // 	if err != nil {
// // 		return errors.WithMessage(err, "bar")
// // 	}
// // 	return nil
// // }
//
// func Bar() error {
// 	err := Foo()
// 	return errors.WithMessage(err, "bar")
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		fmt.Printf("err: %+v\n", err)
// 	}
// }

// NOTE: Sentinel error 处理

func Foo() error {
	return io.EOF
}

func Bar() error {
	err := Foo()
	return errors.WithMessage(err, "bar")
}

func main() {
	err := Bar()
	if err != nil {
		if errors.Cause(err) == io.EOF {
			fmt.Println("EOF err")
			return
		}
		fmt.Printf("err: %+v\n", err)
	}
	return
}
