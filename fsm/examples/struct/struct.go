package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/looplab/fsm"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			// NOTE: closed => open
			// 在 open 事件发生之前触发（这里的 open 是指代 open event）
			"before_open": func(_ context.Context, e *fsm.Event) {
				color.Magenta("| before open\t | %s | %s |", e.Src, e.Dst)
			},
			// 任一事件发生之前触发
			"before_event": func(_ context.Context, e *fsm.Event) {
				color.HiMagenta("| before event\t | %s | %s |", e.Src, e.Dst)
			},
			// 在离开 closed 状态时触发
			"leave_closed": func(_ context.Context, e *fsm.Event) {
				color.Cyan("| leave closed\t | %s | %s |", e.Src, e.Dst)
			},
			// 离开任一状态时触发
			"leave_state": func(_ context.Context, e *fsm.Event) {
				color.HiCyan("| leave state\t | %s | %s |", e.Src, e.Dst)
			},
			// 在进入 open 状态时触发（这里的 open 是指代 open state）
			"enter_open": func(_ context.Context, e *fsm.Event) {
				color.Green("| enter open\t | %s | %s |", e.Src, e.Dst)
			},
			// 进入任一状态时触发
			"enter_state": func(_ context.Context, e *fsm.Event) {
				color.HiGreen("| enter state\t | %s | %s |", e.Src, e.Dst)
			},
			// 在 open 事件发生之后触发（这里的 open 是指代 open event）
			"after_open": func(_ context.Context, e *fsm.Event) {
				color.Yellow("| after open\t | %s | %s |", e.Src, e.Dst)
			},
			// 任一事件结束后触发
			"after_event": func(_ context.Context, e *fsm.Event) {
				color.HiYellow("| after event\t | %s | %s |", e.Src, e.Dst)
			},

			/*
				// NOTE: open => closed
				// 事件触发之前 -> 事件
				"before_close": func(_ context.Context, e *fsm.Event) {
					color.Magenta("| before close\t | %s | %s |", e.Src, e.Dst)
				},
				// 离开当前状态
				"leave_open": func(_ context.Context, e *fsm.Event) {
					color.Cyan("| leave open\t | %s | %s |", e.Src, e.Dst)
				},
				// 进入新状态 -> 状态
				"enter_closed": func(_ context.Context, e *fsm.Event) {
					color.Green("| enter closed\t | %s | %s |", e.Src, e.Dst)
				},
				// 事件完成以后 -> 事件
				"after_close": func(_ context.Context, e *fsm.Event) {
					color.Yellow("| after close\t | %s | %s |", e.Src, e.Dst)
				},

				// NOTE: <NEW_STATE> 格式简写
				// 等价于 enter_open
				"open": func(_ context.Context, e *fsm.Event) {
					color.Green("| enter open\t | %s | %s |", e.Src, e.Dst)
				},
				// 等价于 enter_closed
				"closed": func(_ context.Context, e *fsm.Event) {
					color.Green("| enter closed\t | %s | %s |", e.Src, e.Dst)
				},

				// NOTE: <EVENT> 格式简写
				// 等价于 after_close
				"close": func(_ context.Context, e *fsm.Event) {
					color.Yellow("| after close\t | %s | %s |", e.Src, e.Dst)
				},

				// NOTE: 定义一个未知事件（无效）
				"unknown": func(_ context.Context, e *fsm.Event) {
					color.Red("unknown event\t | %s | %s |", e.Src, e.Dst)
				},
			*/
		},
	)

	return d
}

/*
func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}
*/

func main() {
	door := NewDoor("heaven")

	color.White("--------- closed to open ---------")
	color.White("| event\t\t | src\t  | dst\t |")
	color.White("----------------------------------")

	err := door.FSM.Event(context.Background(), "open")
	if err != nil {
		fmt.Println(err)
	}
	color.White("----------------------------------")

	/*
		color.White("--------- open to closed ---------")
		color.White("| event\t\t | src\t  | dst\t |")
		color.White("----------------------------------")

		err = door.FSM.Event(context.Background(), "close")
		if err != nil {
			fmt.Println(err)
		}
		color.White("----------------------------------")

		err = door.FSM.Event(context.Background(), "unknown")
		if err != nil {
			fmt.Println(err)
		}
	*/
}
