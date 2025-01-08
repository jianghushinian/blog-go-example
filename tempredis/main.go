package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/stvp/tempredis"
)

func main() {
	// 创建并启动一个临时的 Redis 实例
	server, err := tempredis.Start(tempredis.Config{
		"port": "0", // 自动分配端口
	})
	if err != nil {
		log.Fatalf("Failed to start tempredis: %v", err)
	}

	// 放在 defer 中执行，避免阻塞 main goroutine
	defer func() {
		fmt.Println("====================== stdout ======================")
		fmt.Println(server.Stdout())
		fmt.Println("====================== stderr ======================")
		fmt.Println(server.Stderr())
	}()

	// main 退出时关闭 redis-server
	defer server.Term()

	// 获取 Redis 的地址
	fmt.Println("Redis server is running at", server.Socket())

	// 连接临时的 Redis 实例
	client := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	// 使用 Redis 实例
	client.Set("name", "jianghushinian", time.Second)
	val, err := client.Get("name").Result()
	if err != nil {
		fmt.Println("Get redis key error:", err)
		return
	}
	fmt.Println("name:", val)

	time.Sleep(time.Second)

	// 1s 后 name 已经过期
	val, err = client.Get("name").Result()
	if err != nil {
		fmt.Println("Get redis key error:", err)
		return
	}
	fmt.Println("name:", val)
}
