package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var (
		maxWorkers = runtime.GOMAXPROCS(0)           // 获取系统可用的最大 CPU 核心数
		out        = make([]int, 32)                 // 存储 Collatz 结果
		wg         sync.WaitGroup                    // 用于等待 goroutine 完成
		sem        = make(chan struct{}, maxWorkers) // 用于限制最大并发数
	)

	for i := range out {
		// 通过 sem 管理并发，确保最多只有 maxWorkers 个 goroutine 同时执行
		sem <- struct{}{} // 如果 sem 已满，这里会阻塞，直到有空闲槽位

		// 增加 WaitGroup 计数
		wg.Add(1)

		go func(i int) {
			defer wg.Done()          // goroutine 完成时，减少 WaitGroup 计数
			defer func() { <-sem }() // goroutine 完成时，从 sem 中释放一个槽位

			// 执行 Collatz 步骤计算
			out[i] = collatzSteps(i + 1)
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 输出结果
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
