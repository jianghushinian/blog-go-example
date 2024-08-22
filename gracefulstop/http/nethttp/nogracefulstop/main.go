package main

import (
	"log"
	"net/http"
	"time"
)

// NOTE: 未使用优雅退出的普通 HTTP Server

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

	if err := srv.ListenAndServe(); err != nil {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Stopped serving new connections")
}
