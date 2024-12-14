package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
	}
	var err error

	var wg sync.WaitGroup // 零值可用，不必显式初始化

	for _, url := range urls {
		wg.Add(1) // 增加 WaitGroup 计数器

		// 启动一个 goroutine 来获取 URL
		go func() {
			defer wg.Done() // 当 goroutine 完成时递减 WaitGroup 计数器

			resp, e := http.Get(url)
			// FIXME: 这里存在 race condition，可以加锁或者使用 channel 解决，`./race` 目录下有两个 demo 实现
			if e != nil { // 发生错误返回，并记录该错误
				err = e
				return
			}
			defer resp.Body.Close()
			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
		}()
	}

	// 等待所有 goroutine 执行完成
	wg.Wait()
	if err != nil { // err 会记录最后一个错误
		fmt.Printf("Error: %s\n", err)
	}
}

// $ go run waitgroup/main.go
// fetch url http://www.google.com/ status 200 OK
// fetch url http://www.golang.org/ status 200 OK
// Error: Get "http://www.somestupidname.com/": dial tcp: lookup www.somestupidname.com: no such host
