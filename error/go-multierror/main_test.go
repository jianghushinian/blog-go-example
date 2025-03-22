package main

import (
	"errors"
	"sync"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	errs := &multierror.Error{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs = multierror.Append(errs, errors.New("error"))
		}()
	}

	wg.Wait()
	// 预期 100 个错误，实际输出可能 < 100
	assert.Equal(t, 100, len(errs.Errors)) // 测试失败
}

func TestConcurrencyWithChannel(t *testing.T) {
	var wg sync.WaitGroup
	errCh := make(chan error, 100)
	errs := &multierror.Error{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- errors.New("error") // 通过 channel 安全发送 err
		}()
	}

	// 开启子 goroutine 等待并发程序执行完成
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// main goroutine 从 channel 收到 err 并完成聚合
	for err := range errCh {
		errs = multierror.Append(errs, err) // 单 goroutine 聚合，无竞争
	}

	// 预期 100 个错误，实际输出也是 100
	assert.Equal(t, 100, len(errs.Errors)) // 测试通过
}
