package main

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

// func Foo() error {
// 	return nil
// }
//
// func main() {
// 	err := Foo()
// 	// slog.Info(err.Error())
// 	fmt.Printf("INFO: call foo: %s\n", err)
// }

// func Foo() error {
// 	return errors.New("foo error")
// }
//
// func Bar() error {
// 	return Foo()
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		slog.Error(err.Error())
// 	}
// }

// func Foo() error {
// 	return errors.New("foo error")
// }
//
// func Bar() error {
// 	err := Foo()
// 	if err != nil {
// 		// NOTE: 服务降级，记录日志
// 		slog.Error(err.Error())
// 		return nil
// 	}
// 	// do something
// 	return nil
// }
//
// func main() {
// 	err := Bar()
// 	if err != nil {
// 		slog.Error(err.Error())
// 	}
// }

func Foo() error {
	return errors.New("foo error")
}

func Bar() error {
	err := Foo()
	return errors.WithMessage(err, "Bar")
}

func main() {
	err := Bar()
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}
