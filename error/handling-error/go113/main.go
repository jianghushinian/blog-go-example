package main

import (
	"errors"
	"fmt"
	"io"
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
// 		return fmt.Errorf("bar: %w", err)
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

// NOTE: Sentinel error 处理
//
// func Foo() error {
// 	return io.EOF
// }
//
// func Bar() error {
// 	err := Foo()
// 	if err != nil {
// 		return fmt.Errorf("bar: %w", err)
// 	}
// 	return nil
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		if errors.Unwrap(err) == io.EOF {
// 			fmt.Println("EOF err")
// 			return
// 		}
// 		fmt.Printf("err: %+v\n", err)
// 	}
// 	return
// }

// NOTE: errors.Is
//
// func Foo() error {
// 	return io.EOF
// }
//
// func Bar() error {
// 	err := Foo()
// 	if err != nil {
// 		return fmt.Errorf("bar: %w", err)
// 	}
// 	return nil
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		// if err == io.EOF {
// 		if errors.Is(err, io.EOF) {
// 			fmt.Println("EOF err")
// 			return
// 		}
// 		fmt.Printf("err: %+v\n", err)
// 	}
// 	return
// }

// NOTE: errors.As

type MyError struct {
	msg string
	err error
}

func (e *MyError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func Foo() error {
	return &MyError{
		msg: "foo",
		err: io.EOF,
	}
}

func Bar() error {
	err := Foo()
	if err != nil {
		return fmt.Errorf("bar: %w", err)
	}
	return nil
}

func main() {
	err := Bar()
	if err != nil {
		var e *MyError
		if errors.As(err, &e) {
			fmt.Printf("EOF err: %s\n", e)
			return
		}
		fmt.Printf("err: %+v\n", err)
	}
	return
}
