package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sync/semaphore"
)

// ref: https://pkg.go.dev/golang.org/x/sync@v0.10.0/semaphore#example-package-WorkerPool
// 这是 semaphore 包官方的“worker pool”模式示例
// 演示如何使用信号量来限制并行任务中运行的 goroutine 数量

func main() {
	ctx := context.TODO()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)                    // worker pool 支持的最大 worker 数量，取当前机器 CPU 核心数
		sem        = semaphore.NewWeighted(int64(maxWorkers)) // 信号量，资源总数即为最大 worker 数量
		out        = make([]int, 32)                          // 总任务数量
	)

	// 一次最多启动 maxWorkers 数量个 goroutine 计算输出
	for i := range out {
		// 当最大工作数 maxWorkers 个 goroutine 正在执行时，Acquire 会阻塞直到其中一个 goroutine 完成
		if err := sem.Acquire(ctx, 1); err != nil { // 请求资源
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		// 开启新的 goroutine 执行计算任务
		go func(i int) {
			defer sem.Release(1)         // 任务执行完成后释放资源
			out[i] = collatzSteps(i + 1) // 执行 Collatz 步骤计算
		}(i)
	}

	// 获取所有的 tokens 以等待全部 goroutine 执行完成
	// 如果已经通过其他方式（例如 errgroup.Group）在等待工作线程，可以省略这最后一次 Acquire 调用
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	fmt.Println(out)

}

// collatzSteps computes the number of steps to reach 1 under the Collatz
// conjecture. (See https://en.wikipedia.org/wiki/Collatz_conjecture.)
func collatzSteps(n int) (steps int) {
	if n <= 0 {
		panic("nonpositive input")
	}

	for ; n > 1; steps++ {
		if steps < 0 {
			panic("too many steps")
		}

		if n%2 == 0 {
			n /= 2
			continue
		}

		const maxInt = int(^uint(0) >> 1)
		if n > (maxInt-1)/3 {
			panic("overflow")
		}
		n = 3*n + 1
	}

	return steps
}
