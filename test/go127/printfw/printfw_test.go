// Package printfw 演示 Go 1.27 printf vet 检查的新增规则：
// fmt.Errorf 的 %w 操作数不应是"实现了 error 的类型的指针"（*E，其中 E 实现 error）。
//
// 为什么这是 bug：%w 包上指针后，错误链里的动态类型是 *myErr 而不是 myErr。
// 于是 errors.Is(err, myErr{...})、类型断言 err.(myErr)、errors.As 取值类型
// 全部匹配失败，而且编译期毫无提示。最典型的惨案是 fmt.Errorf("%w", &err)
// （err 是 error 变量），把 *error 塞进错误链，errors.Is/As 整条路报废。
//
// printf 检查本来就在 go test 默认 vet 列表里，这个子检查只对 go >= 1.27
// 的文件生效（按 go.mod 的 go 指令判断）：go.mod 升到 1.27 它才拦你，
// go.mod 停在老版本的项目不受影响。
//
// 验证：
//
//	go test .                      # 通过
//	go test -tags printfw_demo .   # 挂：vet 报 %w wants operand of error type myErr...
package printfw

import (
	"errors"
	"fmt"
	"testing"
)

type myErr struct {
	msg string
}

// Error 使用值接收者：myErr 和 *myErr 都实现 error。
func (e myErr) Error() string { return e.msg }

// TestCorrectWrapping 演示正确写法：包装值本身，动态类型就是 myErr，
// errors.Is/As 按预期工作。
func TestCorrectWrapping(t *testing.T) {
	err := fmt.Errorf("operation failed: %w", myErr{msg: "boom"})
	if !errors.Is(err, myErr{msg: "boom"}) {
		t.Fatalf("errors.Is should match the value type, got: %v", err)
	}
	var target myErr
	if !errors.As(err, &target) {
		t.Fatalf("errors.As should match the value type, got: %v", err)
	}
}

// TestPointerWrappingDefeatsErrorsIs 演示 %w 里放指针后 errors.Is/As 匹配不上。
func TestPointerWrappingDefeatsErrorsIs(t *testing.T) {
	// 直接写 fmt.Errorf("%w", &myErr{}) 会被 vet 拦下；
	// 改走 error 接口变量，静态类型是 error，vet 便不会报错，但检查结果依然不匹配。
	var wrapped error = &myErr{msg: "boom"}
	err := fmt.Errorf("operation failed: %w", wrapped)

	// 链里存的是 *myErr，按 myErr 找永远找不到
	if errors.Is(err, myErr{msg: "boom"}) {
		t.Fatal("expected match to fail: the error chain holds *myErr, not myErr")
	}
	var v myErr
	if errors.As(err, &v) {
		t.Fatal("expected errors.As with the value type to fail")
	}

	// 只有按 *myErr 找才找得到
	var pv *myErr
	if !errors.As(err, &pv) {
		t.Fatal("expected errors.As with the pointer type to succeed")
	}
}
