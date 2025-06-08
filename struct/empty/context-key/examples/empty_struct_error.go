//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
)

// NOTE: 空结构体作为 context key 的错误用法

func main() {
	ctx := context.Background()

	key1 := struct{}{}
	ctx = context.WithValue(ctx, key1, "data1")
	fmt.Printf("key1 data: %s\n", ctx.Value(key1))

	key2 := struct{}{}
	ctx = context.WithValue(ctx, key2, "data2")
	fmt.Printf("key2 data: %s\n", ctx.Value(key2))

	// 再次查看 key1 对应的 value
	fmt.Printf("key1 data: %s\n", ctx.Value(key1))
}
