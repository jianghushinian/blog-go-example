package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/looplab/fsm"
)

// NOTE: 异步状态转换

func main() {
	// 构造有限状态机
	f := fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		fsm.Callbacks{
			// 注册 leave_<OLD_STATE> 回调函数
			"leave_start": func(_ context.Context, e *fsm.Event) {
				e.Async() // NOTE: 标记为异步，触发事件时不进行状态转换
			},
		},
	)

	// NOTE: 触发 run 事件，但不会完整状态转换
	err := f.Event(context.Background(), "run")

	// NOTE: Sentinel Error `fsm.AsyncError` 标识异步状态转换
	var asyncError fsm.AsyncError
	ok := errors.As(err, &asyncError)
	if !ok {
		panic(fmt.Sprintf("expected error to be 'AsyncError', got %v", err))
	}

	// NOTE: 主动执行状态转换操作
	if err = f.Transition(); err != nil {
		panic(fmt.Sprintf("Error encountered when transitioning: %v", err))
	}

	// NOTE: 当前状态
	fmt.Printf("current state: %s\n", f.Current())
}
