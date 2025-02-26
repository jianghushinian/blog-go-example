package miniredislock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLockAlreadyExpired = errors.New("miniredislock: failed to unlock, lock was already expired")

// MiniRedisMutex 一个微型的 Redis 分布式锁
type MiniRedisMutex struct {
	name   string        // 会作为分布式锁在 Redis 中的 key
	expiry time.Duration // 锁过期时间
	conn   redis.Cmdable // Redis Client
}

// NewMutex 创建 Redis 分布式锁
func NewMutex(name string, expiry time.Duration, conn redis.Cmdable) *MiniRedisMutex {
	return &MiniRedisMutex{name, expiry, conn}
}

// Lock 加锁
func (m *MiniRedisMutex) Lock(ctx context.Context, value string) (bool, error) {
	reply, err := m.conn.SetNX(ctx, m.name, value, m.expiry).Result()
	if err != nil {
		return false, err
	}
	return reply, nil
}

// 释放锁的 lua 脚本，保证并发安全
var deleteScript = `
	local val = redis.call("GET", KEYS[1])
	if val == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	elseif val == false then
		return -1
	else
		return 0
	end
`

// Unlock 释放锁
func (m *MiniRedisMutex) Unlock(ctx context.Context, value string) (bool, error) {
	// 执行 lua 脚本，Redis 会保证其并发安全
	status, err := m.conn.Eval(ctx, deleteScript, []string{m.name}, value).Result()
	if err != nil {
		return false, err
	}
	if status == int64(-1) {
		return false, ErrLockAlreadyExpired
	}
	return status != int64(0), nil
}
