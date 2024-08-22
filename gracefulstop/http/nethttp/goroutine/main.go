package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// NOTE: 当 HTTP Handler 中存在调用 goroutine 时的优雅退出

type Service struct {
	wg sync.WaitGroup
}

func (s *Service) FakeSendEmail() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered panic: %v\n", err)
			}
		}()

		log.Println("Goroutine enter")
		time.Sleep(time.Second * 5)
		log.Println("Goroutine exit")
	}()
}

func (s *Service) GracefulStop(ctx context.Context) {
	log.Println("Waiting for service to finish")
	quit := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(quit)
	}()
	select {
	case <-ctx.Done():
		log.Println("context was marked as done earlier, than user service has stopped")
	case <-quit:
		log.Println("Service finished")
	}
}

// Handler curl "http://localhost:8000/sleep?duration=5s"
func (s *Service) Handler(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration(r.FormValue("duration"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// svc.FakeSendEmail()

	time.Sleep(duration)

	// 模拟需要异步执行的代码，比如注册接口异步发送邮件、发送 Kafka 消息等
	s.FakeSendEmail()

	_, _ = w.Write([]byte("Welcome HTTP Server"))
}

func main() {
	srv := &http.Server{
		Addr: ":8000",
	}

	svc := &Service{}
	http.HandleFunc("/sleep", svc.Handler)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// 错误写法
	// srv.RegisterOnShutdown(func() {
	// 	svc.GracefulStop(ctx)
	// })

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

	// 注意：保险起见，svc.GracefulStop(ctx) 调用不能放在 srv.RegisterOnShutdown 中注册执行
	// 因为 svc.Handler 中执行到 time.Sleep(duration) 时，还没开始执行 svc.FakeSendEmail()
	// 这时如果按 `Ctrl + C` 退出程序，srv.Shutdown(ctx) 内部会先执行 srv.RegisterOnShutdown
	// 注册的函数，svc.GracefulStop 会立即执行完成并退出，之后等待几秒 svc.Handler 中的逻辑才会
	// 走到 svc.FakeSendEmail()，此时已经无法实现优雅退出 goroutine 了
	svc.GracefulStop(ctx)
	log.Println("HTTP server graceful shutdown completed")
}
