package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

func main() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
	}

	// 创建一个带有 context 的 errgroup
	// 任何一个 goroutine 返回非 nil 的错误，或 Wait() 等待所有 goroutine 完成后，context 都会被取消
	g, ctx := errgroup.WithContext(context.Background())

	// 创建一个 map 来保存结果
	var result sync.Map

	for _, url := range urls {
		// 根据实际情况调整，经测试 sleep 1s 情况下，www.google.com 可以返回，www.golang.org 来不及返回
		// if strings.Contains(url, "somestupidname.com") {
		// 	time.Sleep(time.Second)
		// }

		// 使用 errgroup 启动一个 goroutine 来获取 URL
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err // 发生错误，返回该错误
			}

			// 发起请求
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err // 发生错误，返回该错误
			}
			defer resp.Body.Close()

			// 保存每个 URL 的响应状态码
			result.Store(url, resp.Status)
			return nil // 返回 nil 表示成功
		})
	}

	// 等待所有 goroutine 完成并返回第一个错误（如果有）
	if err := g.Wait(); err != nil {
		fmt.Println("Error: ", err)
	}

	// 所有 goroutine 都执行完成，遍历并打印成功的结果
	result.Range(func(key, value any) bool {
		fmt.Printf("fetch url %s status %s\n", key, value)
		return true
	})
}

// $ go run examples/withcontext/main.go
// Error:  Get "http://www.somestupidname.com/": dial tcp: lookup www.somestupidname.com: no such host
// fetch url http://www.google.com/ status 200 OK
