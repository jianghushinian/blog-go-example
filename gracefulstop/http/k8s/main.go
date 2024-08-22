package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	// genericapiserver "github.com/jianghushinian/blog-go-example/gracefulstop/pkg/server"
)

// NOTE: K8s 优雅退出方案

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
		_, _ = w.Write([]byte("Welcome HTTP Server"))
	})

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// NOTE: 只需要替换这 3 行代码，Gin 版本同理
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// 可以直接丢弃，context.Context.Done() 返回的就是普通空结构体
	<-genericapiserver.SetupSignalHandler()

	// 另一个玩法
	// c := genericapiserver.SetupSignalContext()
	// <-c.Done()
	// log.Println(c.Err()) //  context canceled
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We received an SIGINT/SIGTERM signal, shut down.
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server graceful shutdown completed")
}
