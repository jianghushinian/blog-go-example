package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("main enter")

	// time.Sleep(time.Second)

	quit := make(chan os.Signal, 1)
	// 注册需要关注的信号：SIGINT、SIGTERM、SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞当前 goroutine 等待信号
	sig := <-quit
	fmt.Printf("received signal: %d-%s\n", sig, sig)

	fmt.Println("main exit")
}
