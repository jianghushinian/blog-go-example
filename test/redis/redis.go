package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	smsCaptchaExpire    = 5 * time.Minute
	smsCaptchaKeyPrefix = "sms:captcha:%s"

	authTokenExpire    = 24 * time.Hour
	authTokenKeyPrefix = "auth:token:%s"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SetSmsCaptchaToRedis(ctx context.Context, redis *redis.Client, mobile, captcha string) error {
	key := fmt.Sprintf(smsCaptchaKeyPrefix, mobile)
	return redis.Set(ctx, key, captcha, smsCaptchaExpire).Err()
}

func GetSmsCaptchaFromRedis(ctx context.Context, redis *redis.Client, mobile string) (string, error) {
	key := fmt.Sprintf(smsCaptchaKeyPrefix, mobile)
	return redis.Get(ctx, key).Result()
}

func SetAuthTokenToRedis(ctx context.Context, redis *redis.Client, token, mobile string) error {
	key := fmt.Sprintf(authTokenKeyPrefix, token)
	return redis.Set(ctx, key, mobile, authTokenExpire).Err()
}

func GetAuthTokenFromRedis(ctx context.Context, redis *redis.Client, token string) (string, error) {
	key := fmt.Sprintf(authTokenKeyPrefix, token)
	return redis.Get(ctx, key).Result()
}

func GenerateToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token)[:length], nil
}

func Login(mobile, smsCode string, rdb *redis.Client, generateToken func(int) (string, error)) (string, error) {
	ctx := context.Background()

	// 查找验证码
	captcha, err := GetSmsCaptchaFromRedis(ctx, rdb, mobile)
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("invalid sms code or expired")
		}
		return "", err
	}

	if captcha != smsCode {
		return "", fmt.Errorf("invalid sms code")
	}

	// 登录，生成 token 并写入 Redis
	token, _ := generateToken(32)
	err = SetAuthTokenToRedis(ctx, rdb, token, mobile)
	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	rdb := NewRedisClient()
	token, err := Login("13800001111", "123456", rdb, GenerateToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}
