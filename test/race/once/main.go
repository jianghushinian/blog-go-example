package main

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	id int
}

var instance *Client

func GetClient() *Client {
	if instance == nil {
		time.Sleep(10 * time.Millisecond)
		instance = &Client{id: 1}
	}
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
