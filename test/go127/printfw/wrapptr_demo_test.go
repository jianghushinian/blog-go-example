//go:build printfw_demo

package printfw

import (
	"fmt"
	"testing"
)

// TestWrapPointer 触发 Go 1.27 新增的 printf 子检查：
//
//	%w wants operand of error type myErr, not pointer type *myErr (defeats errors.Is)
//
// 默认（不加 -tags printfw_demo）本文件被构建约束排除，
// 所以日常 go test . 保持绿色。
func TestWrapPointer(t *testing.T) {
	err := fmt.Errorf("operation failed: %w", &myErr{msg: "boom"})
	t.Logf("wrapped: %v", err)
}
