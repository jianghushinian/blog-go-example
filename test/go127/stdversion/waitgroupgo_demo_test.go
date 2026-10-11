//go:build stdversion_demo

package stdversion

import (
	"sync"
	"testing"
)

// TestWaitGroupGoTooNew 使用 sync.WaitGroup.Go——Go 1.25 才加入的 API。
// 本模块 go.mod 声明 go 1.24，vet stdversion 报：
//
//	sync.(*WaitGroup).Go requires go1.25 or later (module is go1.24)
//
// 默认（不加 -tags stdversion_demo）本文件被构建约束排除，
// 所以日常 go test . 保持绿色。
func TestWaitGroupGoTooNew(t *testing.T) {
	var wg sync.WaitGroup
	wg.Go(func() {})
	wg.Wait()
}
