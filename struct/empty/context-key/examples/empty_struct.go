//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
)

// NOTE: 使用空结构体作为 context key
type emptyKey struct{}
type anotherEmpty struct{}

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, emptyKey{}, "empty struct data")

	fmt.Printf("empty data: %s\n", ctx.Value(emptyKey{}))

	ctx = context.WithValue(ctx, anotherEmpty{}, "another empty struct data")

	fmt.Printf("another empty data: %s\n", ctx.Value(anotherEmpty{}))

	// 再次查看 emptyKey 对应的 value
	fmt.Printf("empty data: %s\n", ctx.Value(emptyKey{}))
}
