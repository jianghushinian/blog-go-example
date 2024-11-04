package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// 创建一个 errgroup.Group
	var g errgroup.Group
	// 设置最大并发限制为 3
	g.SetLimit(3)

	// 启动 10 个 goroutine
	for i := 1; i <= 10; i++ {
		if g.TryGo(func() error {
			// 打印正在运行的 goroutine
			fmt.Printf("Goroutine %d is starting\n", i)
			time.Sleep(2 * time.Second) // 模拟工作
			fmt.Printf("Goroutine %d is done\n", i)
			return nil
		}) {
			// 如果成功启动，打印提示
			fmt.Printf("Goroutine %d started successfully\n", i)
		} else {
			// 如果达到并发限制，打印提示
			fmt.Printf("Goroutine %d could not start (limit reached)\n", i)
		}
	}

	// 等待所有 goroutine 完成
	if err := g.Wait(); err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
	}

	fmt.Println("All goroutines complete.")
}

// $ go run examples/trygo/main.go
// Goroutine 1 started successfully
// Goroutine 1 is starting
// Goroutine 2 is starting
// Goroutine 2 started successfully
// Goroutine 3 started successfully
// Goroutine 4 could not start (limit reached)
// Goroutine 5 could not start (limit reached)
// Goroutine 6 could not start (limit reached)
// Goroutine 7 could not start (limit reached)
// Goroutine 8 could not start (limit reached)
// Goroutine 9 could not start (limit reached)
// Goroutine 10 could not start (limit reached)
// Goroutine 3 is starting
// Goroutine 2 is done
// Goroutine 3 is done
// Goroutine 1 is done
// All goroutines complete.
