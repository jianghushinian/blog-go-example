package main

import "fmt"

func Abs(x int) int {
	fmt.Printf(">>> call abs(%d)\n", x)
	if x < 0 {
		return -x
	}
	return x
}
