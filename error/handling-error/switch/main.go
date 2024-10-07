package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("example.txt")
	if err != nil {
		return
	}

	b := bufio.NewReader(f)

	data, err := b.Peek(10)
	if err != nil {
		switch err {
		case bufio.ErrNegativeCount:
			// do something
			return
		case bufio.ErrBufferFull:
			// do something
			return
		default:
			// do something
			return
		}
	}
	fmt.Println(string(data))
}
