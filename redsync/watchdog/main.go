package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"                  // 引入 redsync 库，用于实现基于 Redis 的分布式锁
	"github.com/go-redsync/redsync/v4/redis/goredis/v9" // 引入 redsync 的 goredis 连接池
	goredislib "github.com/redis/go-redis/v9"           // 引入 go-redis 库，用于与 Redis 服务器通信
)

func main() {
	// 创建一个 Redis 客户端
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "localhost:36379", // Redis 服务器地址
		Password: "nightwatch",
	})

	// 使用 go-redis 客户端创建一个 redsync 连接池
	pool := goredis.NewPool(client)

	// 创建一个 redsync 实例，用于管理分布式锁
	rs := redsync.New(pool)

	// 创建一个名为 "test-redsync" 的互斥锁（Mutex）
	mutex := rs.NewMutex("test-redsync", redsync.WithExpiry(5*time.Second))

	// 创建一个上下文（context），一般用于控制锁的超时和取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取锁，如果获取失败（例如锁已被其他进程持有），会返回错误
	if err := mutex.LockContext(ctx); err != nil {
		panic(err) // 如果获取锁失败，程序会 panic
	}

	// 看门狗，实现锁自动续约
	stopCh := make(chan struct{})
	ticker := time.NewTicker(2 * time.Second) // 每隔 2s 续约一次
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				// 续约，延长锁的过期时间
				if ok, err := mutex.ExtendContext(ctx); !ok || err != nil {
					slog.Error("Failed to extend mutex", "err", err, "status", ok)
				} else {
					slog.Info("Successfully extend mutex")
				}
			case <-stopCh:
				slog.Info("Exiting mutex watchdog")
				return
			}
		}
	}()

	// 执行业务逻辑
	time.Sleep(6 * time.Second)

	// 通知看门狗停止自动续期
	stopCh <- struct{}{}

	// 释放锁，如果释放失败（例如锁已过期或不属于当前进程），会返回错误
	if _, err := mutex.UnlockContext(ctx); err != nil {
		panic(err) // 如果释放锁失败，程序会 panic
	}
}
