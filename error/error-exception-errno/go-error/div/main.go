package main

import (
	"errors"
	"fmt"
)

func div(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	defer func() {
		fmt.Println("release resources")
	}()

	result, err := div(1, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
