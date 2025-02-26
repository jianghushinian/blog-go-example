package main

import (
	"fmt"
	"time"

	goredislib "github.com/redis/go-redis/v9"
	"golang.org/x/net/context"

	"github.com/jianghushinian/blog-go-example/redsync/miniredislock"
)

func main() {
	// 创建一个 Redis 客户端
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "localhost:36379", // Redis 服务器地址
		Password: "nightwatch",
	})
	defer client.Close()

	// 创建一个名为 "test-miniredislock" 的互斥锁
	mutex := miniredislock.NewMutex("test-miniredislock", 5*time.Second, client)

	ctx := context.Background()
	// 互斥锁的值应该是一个随机值
	value := "random-string"

	// 获取锁
	_, err := mutex.Lock(ctx, value)
	if err != nil {
		panic(err)
	}

	// 执行业务逻辑
	fmt.Println("do something...")
	time.Sleep(3 * time.Second)

	// 释放自己持有的锁
	_, err = mutex.Unlock(ctx, value)
	if err != nil {
		panic(err)
	}
}
