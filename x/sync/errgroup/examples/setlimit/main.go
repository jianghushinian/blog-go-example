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
		g.Go(func() error {
			// 打印正在运行的 goroutine
			fmt.Printf("Goroutine %d is starting\n", i)
			time.Sleep(2 * time.Second) // 模拟任务耗时
			fmt.Printf("Goroutine %d is done\n", i)
			return nil
		})
	}

	// 等待所有 goroutine 完成
	if err := g.Wait(); err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
	}

	fmt.Println("All goroutines complete.")
}

// $ go run examples/setlimit/main.go
// Goroutine 3 is starting
// Goroutine 1 is starting
// Goroutine 2 is starting
// Goroutine 2 is done
// Goroutine 1 is done
// Goroutine 5 is starting
// Goroutine 3 is done
// Goroutine 6 is starting
// Goroutine 4 is starting
// Goroutine 6 is done
// Goroutine 5 is done
// Goroutine 8 is starting
// Goroutine 4 is done
// Goroutine 7 is starting
// Goroutine 9 is starting
// Goroutine 9 is done
// Goroutine 8 is done
// Goroutine 10 is starting
// Goroutine 7 is done
// Goroutine 10 is done
// All goroutines complete.
