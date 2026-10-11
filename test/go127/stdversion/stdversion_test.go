// Package stdversion 演示 Go 1.27 起 go test 默认启用的 stdversion vet 检查。
//
// 背景：go.mod 的 go 指令（本模块是 go 1.24）声明了模块的最低 Go 版本。
// 但你本地装的是 go1.27，编译器和标准库都是 1.27 的——随手用了 1.25 才
// 加入的 API（如 sync.WaitGroup.Go）也能编译通过、测试全绿。等到有人用
// go1.24 工具链构建（CI 矩阵、被老版本模块依赖）才会挂，而你本地永远发现不了。
//
// stdversion 检查对照"文件生效的 Go 版本"（go.mod 的 go 指令或文件的
// //go:build go1.X 标签）与"所用标准库符号的引入版本"，符号太新就报错。
// 该分析器 Go 1.26 已存在（可用 go test -vet=stdversion 手动启用），
// Go 1.27 起进入 go test 的默认 vet 列表——这就是"升级 1.27 后老项目
// 测试突然挂"的破坏性变化之一。
//
// 关闭方法：go test -vet=off，或 -vet 列表中去掉 stdversion；
// 根治方法：把 go.mod 的 go 指令升到实际用到的 API 版本。
//
// 验证（go.mod 的 toolchain 行保证用 go1.27 工具链跑）：
//
//	go test .                      # 通过：默认文件只用 go1.24 兼容的写法
//	go test -tags stdversion_demo . # 挂：vet 报 sync.(*WaitGroup).Go requires go1.25 or later
package stdversion

import (
	"sync"
	"testing"
)

// TestGo124CompatibleWrite 是 go 1.24 约束下的正确写法：Add/Done/Wait。
func TestGo124CompatibleWrite(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()
}
