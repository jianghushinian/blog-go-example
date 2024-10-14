package main

import (
	"fmt"
	"sync"
)

// func f() {
// 	// recover()
// 	// defer recover()
// 	defer func() {
// 		recover()
// 	}()
//
// 	defer fmt.Println("defer 1")
// 	fmt.Println(1)
// 	panic("woah")
// 	defer fmt.Println("defer 2")
// 	fmt.Println(2)
// }

// NOTE: 直接调用没啥用，返回 <nil>

// func f() {
// 	d := recover()
// 	fmt.Println(d)
// }

// NOTE: 下面这个例子将会捕获到 panic，并且输出 panic 信息

// func f() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("recover:", r)
// 		}
// 	}()
// 	panic("woah")
// }

// NOTE: 多个 panic 同时发生时候，只会捕获第一个 panic

// func f() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("recover:", r)
// 		}
// 	}()
// 	panic("woah 1")
// 	panic("woah 2")
// }

// NOTE: 不要在 defer 中出现 panic

// func f() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("recover:", r)
// 		}
// 	}()
//
// 	defer func() {
// 		panic("woah 1")
// 	}()
// 	panic("woah 2")
// }

// func f() {
// 	defer func() {
// 		panic("woah 1")
// 	}()
//
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("recover:", r)
// 		}
// 	}()
//
// 	panic("woah 2")
// }

// NOTE: recover 只能捕获当前 goroutine 中的 panic

// func f() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("recover:", r)
// 		}
// 	}()
//
// 	go func() {
// 		panic("woah")
// 	}()
// 	time.Sleep(1 * time.Second)
// }

// NOTE: panic 转换成 error，防止调用下层代码出现 panic

// func g(i int) (number int, err error) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			var ok bool
// 			err, ok = r.(error)
// 			if !ok {
// 				err = fmt.Errorf("f returns err: %v", r)
// 			}
// 		}
// 	}()
//
// 	number, err = f(i)
// 	return number, err
// }
//
// func f(i int) (int, error) {
// 	if i == 0 {
// 		panic("i=0")
// 	}
// 	return i * i, nil
// }

// NOTE: 以下代码在旧版本 Go 中存在问题，panic(nil) 无法被 recover 捕获
// 在 Go 1.21 中被修复 https://go.dev/doc/go1.21#language

// 默认新版本 Go 已经解决
// $ go run main.go
// panic called with nil argument

// 新版本 Go 可以按照如下方式触发
// $ GODEBUG=panicnil=1 go run main.go

// func f() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	panic(nil)
// 	// panic(new(runtime.PanicNilError))
// }

// NOTE: 并发读写 map 触发 panic，无法被 recover 捕获

// Go 1.19 Release Notes 有提到 https://go.dev/doc/go1.19#runtime
// Go 1.19 版本以后默认不会打印详细堆栈，可以使用如下两种方式打印详细堆栈
// GOTRACEBACK=system go run main.go
// GOTRACEBACK=crash go run main.go

// func f() {
// 	m := map[int]struct{}{}
//
// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				fmt.Println("goroutine 1", err)
// 			}
// 		}()
// 		for {
// 			m[1] = struct{}{}
// 		}
// 	}()
//
// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				fmt.Println("goroutine 2", err)
// 			}
// 		}()
// 		for {
// 			m[1] = struct{}{}
// 		}
// 	}()
//
// 	select {}
// }

// 使用 sync.Map 来解决
func f() {
	m := sync.Map{}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("goroutine 1", err)
			}
		}()
		for {
			m.Store(1, struct{}{})
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("goroutine 2", err)
			}
		}()
		for {
			m.Store(1, struct{}{})
		}
	}()

	select {}
}

func main() {
	f()

	// fmt.Println(g(1))
	// fmt.Println(g(0))

	// {
	// 	var examples = []int{
	// 		1,
	// 		2,
	// 		0,
	// 		4,
	// 	}
	//
	// 	for _, ex := range examples {
	// 		fmt.Printf("g(%d): ", ex)
	// 		nums, err := g(ex)
	// 		if err != nil {
	// 			fmt.Println(err, "\n")
	// 			continue
	// 		}
	// 		fmt.Println("result:", nums, "\n")
	// 	}
	// }
}
