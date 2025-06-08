//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
)

// NOTE: 为了避免 key 冲突，我们通常可以为 key 定义一个业务属性的前缀
const (
	userDataKey = "user-data"
	postDataKey = "post-data"
)

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, userDataKey, "user data")

	fmt.Printf("user-data: %s\n", ctx.Value(userDataKey))

	ctx = context.WithValue(ctx, postDataKey, "post data")

	fmt.Printf("post-data: %s\n", ctx.Value(postDataKey))
}
