//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
)

// NOTE: 这个 key 非常容易冲突
const dataKey = "data"

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, dataKey, "some data")

	fmt.Printf("data: %s\n", ctx.Value(dataKey))

	userDataKey := "data" // 与 dataKey 值相同
	ctx = context.WithValue(ctx, userDataKey, "user data")

	fmt.Printf("user data: %s\n", ctx.Value(userDataKey))

	// 再次查看 dataKey 的值
	fmt.Printf("data: %s\n", ctx.Value(dataKey))
}
