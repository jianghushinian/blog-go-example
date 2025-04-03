package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/air-verse/air/runner"
)

func main() {
	fmt.Printf("args[1]: %v\n", os.Args[1])

	// 注册优雅退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 初始化 HTTP 服务器
	server := &http.Server{Addr: ":8080"}

	// 定义路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🚀 Hello Air! (PID: %d)", os.Getpid()) // PID 用于验证热替换
	})

	// 启动服务协程
	go func() {
		fmt.Printf("Server started at http://localhost:8080 (PID: %d)\n", os.Getpid())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// 阻塞等待终止信号
	<-signalChan
	fmt.Println("Server shutting down...")
	server.Close()
}
