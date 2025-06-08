package main

import (
	"context"
	"fmt"
)

const requestIdKey = "request-id"

func main() {
	ctx := context.Background()

	// NOTE: 通过 context 传递 request id 信息
	// 设置值
	ctx = context.WithValue(ctx, requestIdKey, "req-123")
	// 获取值
	fmt.Printf("request-id: %s\n", ctx.Value(requestIdKey))
}
