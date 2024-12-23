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

	var wg sync.WaitGroup
	errs := make([]error, len(urls)) // 使用 slice 收集错误

	for i, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				errs[i] = fmt.Errorf("failed to fetch %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
		}()
	}

	wg.Wait()

	// 处理所有错误
	for i, err := range errs {
		if err != nil {
			fmt.Printf("fetch url %s error: %s\n", urls[i], err)
		}
	}
}
