package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/panic" {
		panic("url is error")
	}
	// 打印请求的路径
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
	// 创建一个日志实例，写到标准输出
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// 自定义 HTTP Server
	server := &http.Server{
		Addr:     ":8080",
		ErrorLog: logger, // 设置日志记录器
	}

	// 注册处理函数
	http.HandleFunc("/", handler)

	// 启动服务器
	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Println("Server failed to start:", err)
	}
}
