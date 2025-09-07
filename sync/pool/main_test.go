package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 定义一个可跟踪的对象类型
type TraceableObj struct {
	ID    int
	Value string
}

// 全局的 sync.Pool 实例
var objPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("New: Creating a brand new TraceableObj")
		return &TraceableObj{}
	},
}

func TestPoolGC(t *testing.T) {
	// 设置 GOMAXPROCS 为 1，使得测试在一个 P 上执行，减少并发变量
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(runtime.NumCPU()) // 测试结束后恢复

	// 创建第一个对象并放入池中
	fmt.Println("Step 1: Create and Put an object into the pool")
	obj1 := objPool.Get().(*TraceableObj)
	obj1.ID = 1
	obj1.Value = "First Object"
	// 设置 finalizer，当对象被 GC 回收时打印信息
	runtime.SetFinalizer(obj1, func(o *TraceableObj) {
		fmt.Printf("Finalizer Called: Object %d ('%s') has been garbage collected.\n", o.ID, o.Value)
	})
	objPool.Put(obj1)
	fmt.Printf("Object %d put back into pool.\n", obj1.ID)

	// 手动触发一次 GC
	fmt.Println("\nStep 2: Trigger First GC")
	runtime.GC()                       // 对象会从 local 移至 victim，不会被真正回收
	time.Sleep(time.Millisecond * 100) // 给 GC 和 finalizer 调度一点时间

	// 尝试从池中获取，很可能从 victim 中拿到旧对象
	fmt.Println("\nStep 3: Get object after first GC")
	obj2 := objPool.Get().(*TraceableObj)
	fmt.Printf("Got object %d ('%s') after first GC.\n", obj2.ID, obj2.Value)
	objPool.Put(obj2) // 再次放回

	// 连续手动触发二次 GC
	fmt.Println("\nStep 4: Trigger Second GC")
	runtime.GC()
	runtime.GC() // 这次 victim 中的对象会被真正回收
	time.Sleep(time.Millisecond * 100)

	// 再次尝试获取对象
	fmt.Println("\nStep 5: Get object after second GC")
	obj3 := objPool.Get().(*TraceableObj) // 很可能触发 New，因为 victim 被清空了
	fmt.Printf("Got object %d ('%s') after second GC.\n", obj3.ID, obj3.Value)
	objPool.Put(obj3)

	// 额外触发一次 GC 来清理可能剩余的对象（例如 obj3）
	fmt.Println("\nStep 6: Trigger one more GC to clean up if necessary")
	runtime.GC()
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Test completed.")
}
