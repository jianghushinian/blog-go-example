package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// NOTE: 使用 map + struct{} 实现 set

// Set 基于空结构体实现 set
type Set map[string]struct{}

// Add 添加元素到 set
func (s Set) Add(element string) {
	s[element] = struct{}{}
}

// Remove 从 set 中移除元素
func (s Set) Remove(element string) {
	delete(s, element)
}

// Contains 检查 set 中是否包含指定元素
func (s Set) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

// Size 返回 set 大小
func (s Set) Size() int {
	return len(s)
}

// String implements fmt.Stringer
func (s Set) String() string {
	format := "("
	for element := range s {
		format += element + " "
	}
	format = strings.TrimRight(format, " ") + ")"
	return format
}

// NOTE: 无操作的方法接收器

type NoOp struct{}

func (n NoOp) Perform() {
	fmt.Println("Performing no operation.")
}

func main() {
	// 空结构体不占用内存空间
	{
		type Empty struct{}

		var s1 struct{}
		s2 := Empty{}
		s3 := struct{}{}

		fmt.Printf("s1 addr: %p, size: %d\n", &s1, unsafe.Sizeof(s1))
		fmt.Printf("s2 addr: %p, size: %d\n", &s2, unsafe.Sizeof(s2))
		fmt.Printf("s3 addr: %p, size: %d\n", &s3, unsafe.Sizeof(s3))
		fmt.Printf("s1 == s2 == s3: %t\n", s1 == s2 && s2 == s3)
	}

	// 嵌套空结构体同样不占用内存空间
	{
		type Empty struct{}

		type MultiEmpty struct {
			A Empty
			B struct{}
		}

		s1 := Empty{}
		s2 := MultiEmpty{}
		fmt.Printf("s1 addr: %p, size: %d\n", &s1, unsafe.Sizeof(s1))
		fmt.Printf("s2 addr: %p, size: %d\n", &s2, unsafe.Sizeof(s2))
	}

	// 空结构体顺序不同，会影响内存对齐
	{
		type A struct {
			x int
			y string
			z struct{}
		}

		type B struct {
			x int
			z struct{}
			y string
		}

		type C struct {
			z struct{}
			x int
			y string
		}

		a := A{}
		b := B{}
		c := C{}
		fmt.Printf("struct a size: %d\n", unsafe.Sizeof(a))
		fmt.Printf("struct b size: %d\n", unsafe.Sizeof(b))
		fmt.Printf("struct c size: %d\n", unsafe.Sizeof(c))
	}

	// 使用空结构体实现 set
	{
		s := make(Set)

		s.Add("one")
		s.Add("two")
		s.Add("three")

		fmt.Printf("set: %s\n", s)
		fmt.Printf("set size: %d\n", s.Size())
		fmt.Printf("set contains 'one': %t\n", s.Contains("one"))
		fmt.Printf("set contains 'onex': %t\n", s.Contains("onex"))

		s.Remove("one")

		fmt.Printf("set: %s\n", s)
		fmt.Printf("set size: %d\n", s.Size())

		fmt.Printf("set two value: %s, size: %d\n", s["two"], unsafe.Sizeof(s["two"]))
	}

	// 也许有人会认为这样也可以实现 set，但其实 any 是会占用空间的
	{
		s := make(map[string]any)
		s["t1"] = nil
		s["t2"] = struct{}{}
		fmt.Printf("set t1 value: %v, size: %d\n", s["t1"], unsafe.Sizeof(s["t1"]))
		fmt.Printf("set t2 value: %v, size: %d\n", s["t2"], unsafe.Sizeof(s["t2"]))
		fmt.Printf("%T %T\n", s["t1"], s["t2"])
	}

	// set 惯用法
	{
		s := map[string]struct{}{
			"one":   {},
			"two":   {},
			"three": {},
		}
		for element := range s {
			fmt.Println(element)
		}
	}

	// array 不占用空间
	{
		var a [1000000]string
		var b [1000000]struct{}

		fmt.Printf("array a size: %d\n", unsafe.Sizeof(a))
		fmt.Printf("array b size: %d\n", unsafe.Sizeof(b))
	}

	// slice 只占用 `header` 的空间
	{
		var a = make([]string, 1000000)
		var b = make([]struct{}, 1000000)
		fmt.Printf("slice a size: %d\n", unsafe.Sizeof(a))
		fmt.Printf("slice b size: %d\n", unsafe.Sizeof(b))
	}

	// 信号
	{
		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second) // 执行一些操作...
			fmt.Printf("goroutine done\n")
			done <- struct{}{} // 发送完成信号
		}()

		fmt.Printf("waiting...\n")
		<-done // 等待完成
		fmt.Printf("main exit\n")

		// Context 接口的 Done() 方法也使用 struct{}{} 作为信号
		_ = context.Context(nil)
		// Done() <-chan struct{}
	}

	// 信号的另一种实现：占位符
	{
		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second) // 执行一些操作...
			fmt.Printf("goroutine done\n")
			close(done) // 不需要发送 struct{}{}，直接 close，发送完成信号
		}()

		fmt.Printf("waiting...\n")
		<-done // 等待完成
		fmt.Printf("main exit\n")
	}

	// 无操作的方法接收器
	{
		NoOp{}.Perform()
	}

	// 作为接口实现
	{
		// `io.Discard` 的主要作用是提供一个“黑洞”设备，任何写入到 `io.Discard` 的数据都会被消耗掉而不会有任何效果。
		// 这类似于 Unix 中的 `/dev/null` 设备。
		_ = io.Discard

		// e.g., 设置日志输出为 `io.Discard`，忽略所有日志
		log.SetOutput(io.Discard)
		// 这条日志不会在任何地方显示
		log.Println("This log will not be shown anywhere")
	}

	// 标识符
	{
		// `sync.Pool` 包含一个 `noCopy` 属性，`noCopy` 是 Go 源码中禁止复制的检测方法
		// `go vet` 命令可以检测出 `sync.Pool` 是否被复制
		_ = sync.Pool{}

		// 自定义 struct 也可以嵌入 `noCopy` 属性来实现禁止复制
		// ref: blog-go-example/struct/empty/nocpoy/main.go
	}
}
