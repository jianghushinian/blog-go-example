package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	key := struct{}{}
	ctx1 := context.Background()
	ctx2 := context.WithValue(ctx1, key, "value2")
	ctx3, cancel3 := context.WithCancel(ctx1)
	defer cancel3()
	ctx4, cancel4 := context.WithCancel(ctx2)
	defer cancel4()
	ctx5, cancel5 := context.WithTimeout(ctx2, 1*time.Second)
	defer cancel5()
	ctx6 := context.WithoutCancel(ctx3)
	ctx7, cancel7 := context.WithCancel(ctx3)
	defer cancel7()
	ctx8 := context.WithValue(ctx5, key, "value8")
	ctx9, cancel9 := context.WithCancel(ctx6)
	defer cancel9()
	ctx10 := context.WithoutCancel(ctx8)

	// only for compile
	_ = ctx1
	_ = ctx2
	_ = ctx3
	_ = ctx4
	_ = ctx5
	_ = ctx6
	_ = ctx7
	_ = ctx8
	_ = ctx9
	_ = ctx10

	// NOTE: 控制链路，控制是从上向下传播
	// 取消 3，7 被级联取消，6 不支持取消，控制链路被打断，9 不会被取消
	// {
	// 	cancel3()                                             // 取消 3
	// 	fmt.Printf("ctx6: %v, %v\n", ctx6.Done(), ctx6.Err()) // <nil>, <nil>
	// 	fmt.Printf("ctx7: %v, %v\n", ctx7.Done(), ctx7.Err()) // 0x14000102070, context canceled
	// 	fmt.Printf("ctx9: %v, %v\n", ctx9.Done(), ctx9.Err()) // 0x14000102150, <nil>
	// 	cancel9()                                             // 取消 9
	// 	fmt.Printf("ctx9: %v, %v\n", ctx9.Done(), ctx9.Err()) // 0x14000102150, context canceled
	// }

	// 取消 7，3 不会被取消
	// {
	// 	cancel7()
	// 	fmt.Printf("ctx7: %v, %v\n", ctx7.Done(), ctx7.Err()) // 0x14000098070, context canceled
	// 	fmt.Printf("ctx3: %v, %v\n", ctx3.Done(), ctx3.Err()) // 0x140000980e0, <nil>
	// }

	// NOTE: 安全传值，查找值是自下而上
	// 给定 key 查找 value
	// {
	// 	fmt.Printf("ctx2: %s\n", ctx2.Value(key))   // value2
	// 	fmt.Printf("ctx4: %s\n", ctx5.Value(key))   // value2
	// 	fmt.Printf("ctx5: %s\n", ctx5.Value(key))   // value2
	// 	fmt.Printf("ctx8: %s\n", ctx8.Value(key))   // value8
	// 	fmt.Printf("ctx10: %s\n", ctx10.Value(key)) // value8
	// }

	// 5 到期自动取消，不影响给定 key 查找 value
	{
		fmt.Printf("ctx5: %v, %v\n", ctx5.Done(), ctx5.Err()) // 0x14000098150, <nil>
		time.Sleep(2 * time.Second)
		fmt.Printf("ctx5: %v, %v\n", ctx5.Done(), ctx5.Err()) // 0x14000098150, context deadline exceeded

		fmt.Printf("ctx2: %s\n", ctx2.Value(key))   // value2
		fmt.Printf("ctx4: %s\n", ctx5.Value(key))   // value2
		fmt.Printf("ctx5: %s\n", ctx5.Value(key))   // value2
		fmt.Printf("ctx8: %s\n", ctx8.Value(key))   // value8
		fmt.Printf("ctx10: %s\n", ctx10.Value(key)) // value8
	}

	// context.AfterFunc demo
	{
		// {
		// 	ctx := context.Background()
		// 	ctx, cancel := context.WithCancel(ctx)
		// 	f := func() {
		// 		fmt.Println("calling f")
		// 	}
		// 	stop := context.AfterFunc(ctx, f)
		// 	_ = stop
		//
		// 	cancel() // context 取消时会执行 f
		// 	time.Sleep(1 * time.Second)
		// }

		// {
		// 	ctx := context.Background()
		// 	ctx, cancel := context.WithCancel(ctx)
		// 	f := func() {
		// 		fmt.Println("calling f")
		// 	}
		// 	stop := context.AfterFunc(ctx, f)
		//
		// 	stop() // 阻止 f 执行
		//
		// 	cancel() // context 取消时不会执行 f
		// 	time.Sleep(1 * time.Second)
		// }

		// {
		// 	ctx := context.Background()
		// 	ctx, cancel := context.WithCancel(ctx)
		// 	f := func() {
		// 		fmt.Println("calling f")
		// 		time.Sleep(1 * time.Second)
		// 		fmt.Println("called f")
		// 	}
		// 	stop := context.AfterFunc(ctx, f)
		//
		// 	cancel() // context 取消时会执行 f
		// 	stop()   // context 已经被取消，无法停止正在执行的 f
		// 	fmt.Println("stopped")
		// 	time.Sleep(1 * time.Second)
		// }
	}
}
