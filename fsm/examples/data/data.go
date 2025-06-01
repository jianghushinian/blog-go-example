package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

// NOTE: 将 FSM 作为生产者消费者使用

func main() {
	fsm := fsm.NewFSM(
		"idle",
		fsm.Events{
			// 生产者
			{Name: "produce", Src: []string{"idle"}, Dst: "idle"},
			// 消费者
			{Name: "consume", Src: []string{"idle"}, Dst: "idle"},
			// 清理数据
			{Name: "remove", Src: []string{"idle"}, Dst: "idle"},
		},
		fsm.Callbacks{
			// 生产者
			"produce": func(_ context.Context, e *fsm.Event) {
				dataValue := "江湖十年"
				e.FSM.SetMetadata("message", dataValue)
				fmt.Printf("produced data: %s\n", dataValue)
			},
			// 消费者
			"consume": func(_ context.Context, e *fsm.Event) {
				data, ok := e.FSM.Metadata("message")
				if ok {
					fmt.Printf("consume data: %s\n", data)
				}
			},
			// 清理数据
			"remove": func(_ context.Context, e *fsm.Event) {
				e.FSM.DeleteMetadata("message")
				if _, ok := e.FSM.Metadata("message"); !ok {
					fmt.Println("removed data")
				}
			},
		},
	)

	fmt.Printf("current state: %s\n", fsm.Current())

	err := fsm.Event(context.Background(), "produce")
	if err != nil {
		fmt.Printf("produce err: %s\n", err)
	}

	fmt.Printf("current state: %s\n", fsm.Current())

	err = fsm.Event(context.Background(), "consume")
	if err != nil {
		fmt.Printf("consume err: %s\n", err)
	}

	fmt.Printf("current state: %s\n", fsm.Current())

	err = fsm.Event(context.Background(), "remove")
	if err != nil {
		fmt.Printf("remove err: %s\n", err)
	}

	fmt.Printf("current state: %s\n", fsm.Current())
}
