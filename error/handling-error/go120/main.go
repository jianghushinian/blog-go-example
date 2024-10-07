package main

import (
	"errors"
	"fmt"
)

func main() {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err := errors.Join(err1, err2)
	fmt.Println("---------")
	fmt.Println(err)
	fmt.Println("---------")
	if errors.Is(err, err1) {
		fmt.Println("err is err1")
	}
	if errors.Is(err, err2) {
		fmt.Println("err is err2")
	}
}
