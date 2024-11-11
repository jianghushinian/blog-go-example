package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	once := sync.OnceValues(func() ([]byte, error) {
		fmt.Println("Reading file once")
		return os.ReadFile("oncevalues/example_test.go")
	})
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			data, err := once()
			if err != nil {
				fmt.Println("error:", err)
			}
			_ = data // Ignore the data for this example
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
