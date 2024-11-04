package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
	}

	// 使用 errgroup 创建一个新的 goroutine 组
	var g errgroup.Group // 零值可用，不必显式初始化

	for _, url := range urls {
		// 使用 errgroup 启动一个 goroutine 来获取 URL
		g.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err // 发生错误，返回该错误
			}
			defer resp.Body.Close()
			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
			return nil // 返回 nil 表示成功
		})
	}

	// 等待所有 goroutine 完成并返回第一个错误（如果有）
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

// $ go run examples/main.go
// fetch url http://www.google.com/ status 200 OK
// fetch url http://www.golang.org/ status 200 OK
// Error: Get "http://www.somestupidname.com/": dial tcp: lookup www.somestupidname.com: no such host
