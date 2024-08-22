package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NOTE: 实现了优雅退出，不过一定要等请求完成或者超时退出，连续按多次 Ctrl+C 不生效

func main() {
	srv := &http.Server{
		Addr: ":8000",
	}

	// curl "http://localhost:8000/sleep?duration=5s"
	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		time.Sleep(duration)

		// 模拟需要异步执行的代码，比如注册接口异步发送邮件、发送 Kafka 消息等
		// go func() {
		// 	log.Println("Goroutine enter")
		// 	time.Sleep(time.Second * 5)
		// 	log.Println("Goroutine exit")
		// }()

		_, _ = w.Write([]byte("Welcome HTTP Server"))
	})

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// 可以注册一些 hook 函数，比如从注册中心下线逻辑
	srv.RegisterOnShutdown(func() {
		log.Println("Register Shutdown 1")
	})
	srv.RegisterOnShutdown(func() {
		log.Println("Register Shutdown 2")
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We received an SIGINT/SIGTERM/SIGQUIT signal, shut down.
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server graceful shutdown completed")
}
