//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
)

// NOTE: 基于 string 自定义类型，作为 context key
type key1 string
type key2 string

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, key1(""), "data1")
	fmt.Printf("key1 data: %s\n", ctx.Value(key1("")))

	ctx = context.WithValue(ctx, key2(""), "data2")
	fmt.Printf("key2 data: %s\n", ctx.Value(key2("")))

	// 再次查看 key1 对应的 value
	fmt.Printf("key1 data: %s\n", ctx.Value(key1("")))
}
