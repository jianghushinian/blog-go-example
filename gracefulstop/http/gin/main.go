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

	"github.com/gin-gonic/gin"
)

// NOTE: Gin 框架中的优雅退出
// ref: https://gin-gonic.com/zh-cn/docs/examples/graceful-restart-or-stop/

func main() {
	router := gin.Default()

	// curl "http://localhost:8000/sleep?duration=5s"
	router.GET("/sleep", func(c *gin.Context) {
		duration, err := time.ParseDuration(c.Query("duration"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		time.Sleep(duration)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We received an SIGINT/SIGTERM signal, shut down.
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server graceful shutdown completed")
}
