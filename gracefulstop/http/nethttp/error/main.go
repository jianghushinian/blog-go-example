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

// NOTE: 错误写法，不能实现优雅退出效果
// ref: https://pkg.go.dev/net/http@go1.22.0#example-Server.Shutdown
// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed.
// Make sure the program doesn't exit and waits instead for Shutdown to return.

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
		_, _ = w.Write([]byte("Hello World!"))
	})

	go func() {
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
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Stopped serving new connections")
}
