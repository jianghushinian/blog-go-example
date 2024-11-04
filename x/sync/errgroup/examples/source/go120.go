// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 构建约束标识了这个文件是 Go 1.20 版本被加入的
//go:build go1.20

package errgroup

import "context"

// 代理到 context.WithCancelCause
func withCancelCause(parent context.Context) (context.Context, func(error)) {
	return context.WithCancelCause(parent)
}
