package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// 这个变量在 redis.go 文件中读取不到
var rdbClient *redis.Client

func TestLogin_by_container(t *testing.T) {
	// 准备测试数据
	err := SetSmsCaptchaToRedis(context.Background(), rdbClient, "18900001111", "123456")
	assert.NoError(t, err)

	// 测试登录成功情况
	gotToken, err := Login("18900001111", "123456", rdbClient, GenerateToken)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(gotToken))

	// 检查 Redis 中是否存在 token
	gotMobile, err := GetAuthTokenFromRedis(context.Background(), rdbClient, gotToken)
	assert.NoError(t, err)
	assert.Equal(t, "18900001111", gotMobile)
}

func TestLogin_wrong_captcha(t *testing.T) {
	// 准备测试数据
	err := SetSmsCaptchaToRedis(context.Background(), rdbClient, "18900002222", "111111")
	assert.NoError(t, err)

	// 测试验证码不正确情况
	token, err := Login("18900002222", "222222", rdbClient, GenerateToken)
	assert.Equal(t, err, fmt.Errorf("invalid sms code"))
	assert.Equal(t, "", token)
}

// 在容器中运行一个 Redis 服务
func RunWithRedisInContainer() (*redis.Client, func()) {
	ctx := context.Background()

	// 创建容器请求参数
	req := testcontainers.ContainerRequest{
		Image:        "redis:6.0.20-alpine",                      // 指定容器镜像
		ExposedPorts: []string{"6379/tcp"},                       // 指定容器暴露端口
		WaitingFor:   wait.ForLog("Ready to accept connections"), // 等待输出容器 Ready 日志
	}

	// 创建 Redis 容器
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to start container: %s", err.Error()))
	}

	// 获取容器中 Redis 连接地址，e.g. localhost:50351
	endpoint, err := redisC.Endpoint(ctx, "") // 如果暴露多个端口，可以指定第二个参数
	if err != nil {
		panic(fmt.Sprintf("failed to get endpoint: %s", err.Error()))
	}

	// 连接容器中的 Redis
	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	// 返回 Redis Client 和 cleanup 函数
	return client, func() {
		if err := redisC.Terminate(ctx); err != nil {
			panic(fmt.Sprintf("failed to terminate container: %s", err.Error()))
		}
	}
}

func TestMain(m *testing.M) {
	client, f := RunWithRedisInContainer()
	defer f()
	rdbClient = client
	m.Run()
}
