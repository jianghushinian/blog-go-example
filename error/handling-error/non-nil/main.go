package main

import (
	"fmt"
)

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

func returnsError() error {
	var p *MyError = nil
	return p // Will always return a non-nil error.
}

func main() {
	err := returnsError()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("success")
}
