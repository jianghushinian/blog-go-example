package main

import (
	"fmt"
	"net/http"
	"sync"
)

// NOTE: 使用 mutex 解决并发问题

// func main() {
// 	var urls = []string{
// 		"http://www.golang.org/",
// 		"http://www.google.com/",
// 		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
// 	}
// 	var err error
// 	var mu sync.Mutex
//
// 	var wg sync.WaitGroup // 零值可用，不必显式初始化
//
// 	for _, url := range urls {
// 		url := url
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			resp, e := http.Get(url)
// 			if e != nil {
// 				mu.Lock()
// 				if e != nil { // use double-check
// 					err = e
// 				}
// 				mu.Unlock()
// 				return
// 			}
// 			defer resp.Body.Close()
// 			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
// 		}()
// 	}
//
// 	// 等待所有 goroutine 执行完成
// 	wg.Wait()
// 	if err != nil { // err 会记录最后一个错误
// 		fmt.Printf("Error: %s\n", err)
// 	}
// }

// NOTE: 使用 channel 解决并发问题

// func main() {
// 	var urls = []string{
// 		"http://www.golang.org/",
// 		"http://www.google.com/",
// 		"http://www.somestupidname.com/", // 这是一个错误的 URL，会导致任务失败
// 	}
//
// 	var wg sync.WaitGroup
// 	errs := make(chan error, len(urls)) // 使用缓冲通道 channel 收集错误
//
// 	for _, url := range urls {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			resp, err := http.Get(url)
// 			if err != nil {
// 				errs <- fmt.Errorf("failed to fetch %s: %v", url, err)
// 				return
// 			}
// 			defer resp.Body.Close()
// 			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
// 		}()
// 	}
//
// 	wg.Wait()
// 	close(errs)
//
// 	// 处理所有错误
// 	for err := range errs {
// 		if err != nil {
// 			fmt.Printf("Error: %s\n", err)
// 		}
// 	}
// }

// NOTE: 使用 slice 解决并发问题

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
			fmt.Printf("fetch url %s Error: %s\n", urls[i], err)
		}
	}
}
