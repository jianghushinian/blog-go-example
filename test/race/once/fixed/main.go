package main

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	id int
}

var (
	once     sync.Once
	instance *Client
)

func GetClient() *Client {
	once.Do(func() {
		time.Sleep(10 * time.Millisecond)
		instance = &Client{id: 1}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Go(func() {
			c := GetClient()
			fmt.Println(c.id)
		})
	}
	wg.Wait()
}
