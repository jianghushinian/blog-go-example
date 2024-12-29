package main

import (
	"fmt"
	"sync"
	"time"
)

// NOTE: 一个能 run 起来但没什么意义的示例

// func main() {
// 	c := sync.NewCond(&sync.Mutex{})
//
// 	go func() { // 启动一个子 goroutine 进行等待
// 		fmt.Println("wait before")
// 		c.L.Lock()
// 		c.Wait() // 阻塞并等待通知
// 		c.L.Unlock()
// 		fmt.Println("wait after")
// 	}()
//
// 	time.Sleep(time.Second)
//
// 	fmt.Println("signal before")
// 	c.Signal() // 通知唤醒一个阻塞的 goroutine
// 	fmt.Println("signal after")
//
// 	time.Sleep(time.Second) // 确保子 goroutine 执行完成再退出
// }

// NOTE: 一个 sync.Cond 正确用法的示例

func main() {
	c := sync.NewCond(&sync.Mutex{})
	condition := false // 定义一个条件变量

	go func() { // 启动一个子 goroutine 进行等待
		fmt.Println("wait before")
		c.L.Lock()
		for !condition { // 通过循环检查条件是否满足
			c.Wait() // 阻塞并等待通知
		}
		fmt.Println("condition met, continue execution")
		c.L.Unlock()
		fmt.Println("wait after")
	}()

	time.Sleep(time.Second)

	fmt.Println("signal before")
	c.L.Lock()
	condition = true // 改变条件变量的状态
	c.L.Unlock()
	c.Signal() // 通知唤醒一个阻塞的 goroutine
	fmt.Println("signal after")

	time.Sleep(time.Second) // 确保子 goroutine 执行完成再退出
}
