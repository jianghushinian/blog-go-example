package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

func main() {
	task()
	fetch()
}

func task() {
	var g errgroup.Group
	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			if i == 3 {
				return errors.New("task 3 failed")
			}
			if i == 5 {
				return errors.New("task 5 failed")
			}

			// 其他任务继续运行
			fmt.Printf("run task %d\n", i)

			return nil // 正常返回 nil 表示成功
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func fetch() {
	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
	}

	// 创建一个 map 来保存结果
	var result sync.Map

	// 启动多个 goroutine，并发处理多个 URL
	for _, url := range urls {
		// NOTE: 注意这里的 url 需要传递给闭包函数，避免闭包共享变量问题，自 Go 1.22 开始无需考虑此问题
		url := url // https://golang.org/doc/faq#closures_and_goroutines

		// 启动一个 goroutine 来获取 URL
		g.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err // 发生错误，返回该错误
			}
			defer resp.Body.Close()

			// 保存每个 URL 的响应状态码
			result.Store(url, resp.Status)
			return nil
		})
	}

	// 等待所有 goroutine 完成
	if err := g.Wait(); err != nil {
		// 如果有任何一个 goroutine 返回了错误，这里会得到该错误
		fmt.Println("Error: ", err)
	}

	// 所有 goroutine 都执行完成，遍历并打印成功的结果
	result.Range(func(key, value any) bool {
		fmt.Printf("%s: %s\n", key, value)
		return true
	})
}
