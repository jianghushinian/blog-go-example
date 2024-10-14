package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"gorm.io/gorm"
)

// NOTE: 语法

// func f() {
// 	defer fmt.Println("deferred in f")
// 	fmt.Println("calling f")
// }

// NOTE: 多个 defer 执行顺序为 LIFO

// func f() {
// 	defer fmt.Println("deferred in f 1")
// 	defer fmt.Println("deferred in f 2")
// 	defer fmt.Println("deferred in f 3")
// 	fmt.Println("calling f")
// }

// NOTE: 多个 defer 嵌套情况

// func f() {
// 	fmt.Println("1")
//
// 	defer func() {
// 		fmt.Println("2")
// 		defer fmt.Println("3")
// 		fmt.Println("4")
// 	}()
//
// 	fmt.Println("5")
//
// 	defer fmt.Println("6")
//
// 	fmt.Println("7")
// }

// NOTE: defer 改变函数返回值

// f returns 2
// func f() int {
// 	r := 2
// 	defer func() {
// 		fmt.Println("r:", r)
// 		r *= 3
// 	}()
// 	return r
// }

// f returns 6
// func f() (r int) {
// 	r = 2
// 	defer func() {
// 		fmt.Println("r:", r)
// 		r *= 3
// 	}()
// 	return r
// }

// f returns 6
// func f() (r int) {
// 	defer func() {
// 		fmt.Println("r:", r)
// 		r *= 3
// 	}()
// 	return 2
// }

// f returns 2
// func f() (r int) {
// 	defer func(r int) {
// 		fmt.Println("r:", r)
// 		r *= 3
// 	}(r)
// 	return 2
// }

// f returns 2
// func f() (r int) {
// 	x := 2
// 	defer func() {
// 		fmt.Println("r:", r)
// 		fmt.Println("x:", x)
// 		x *= 3
// 	}()
// 	return x
// }

// f returns 6
// func f() (r *int) {
// 	x := 2
// 	defer func() {
// 		fmt.Println("r:", *r)
// 		fmt.Println("x:", x)
// 		x *= 3
// 	}()
// 	return &x
// }

// NOTE: 关闭文件对象

// CopyFile 存在 bug，如果 os.Create 执行失败，函数返回后 src 并没有关闭
// func CopyFile(dstName, srcName string) (written int64, err error) {
// 	src, err := os.Open(srcName)
// 	if err != nil {
// 		return
// 	}
//
// 	dst, err := os.Create(dstName)
// 	if err != nil {
// 		return
// 	}
//
// 	written, err = io.Copy(dst, src)
// 	dst.Close()
// 	src.Close()
// 	return
// }

// CopyFile 丑陋的办法
// func CopyFile(dstName, srcName string) (written int64, err error) {
// 	src, err := os.Open(srcName)
// 	if err != nil {
// 		return
// 	}
//
// 	dst, err := os.Create(dstName)
// 	if err != nil {
// 		src.Close()
// 		return
// 	}
//
// 	written, err = io.Copy(dst, src)
// 	dst.Close()
// 	src.Close()
// 	return
// }

// NOTE: 释放资源 - 关闭文件对象

// CopyFile 正确写法
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// NOTE: 踩坑，一个变量重复赋值

type fakeFile struct {
	name string
}

func (f *fakeFile) Close() error {
	fmt.Println("close:", f)
	return nil
}

// 错误写法：f 变量的值最终是 f2，所以 f2 会被关闭两次，f1 没关闭
func processFile() {
	f := fakeFile{name: "f1"}
	defer f.Close()

	f = fakeFile{name: "f2"}
	defer f.Close()

	fmt.Println("calling processFile")
	return
}

// 解决方案 1
func processFile1() {
	f := fakeFile{name: "f1"}
	defer func(f fakeFile) {
		f.Close()
	}(f)

	f = fakeFile{name: "f2"}
	defer func(f fakeFile) {
		f.Close()
	}(f)

	fmt.Println("calling processFile1")
	return
}

// 解决方案 2
func processFile2() {
	f1 := fakeFile{name: "f1"}
	defer f1.Close()

	f2 := fakeFile{name: "f2"}
	defer f2.Close()

	fmt.Println("calling processFile2")
	return
}

// NOTE: 结构体方法是否使用指针接收者

type User struct {
	name string
}

func (u User) Name() {
	fmt.Println("Name:", u.name)
}

func (u *User) PointName() {
	fmt.Println("PointName:", u.name)
}

// PointName: user2
// Name: user1
// func printUser() {
// 	u := User{name: "user1"}
//
// 	defer u.Name()
// 	defer u.PointName()
//
// 	u.name = "user2"
// }

// Name: user2
// PointName: user2
func printUser() {
	u := User{name: "user1"}

	defer func() {
		u.Name()
		u.PointName()
	}()

	u.name = "user2"
}

// NOTE: defer 执行并不受代码块影响

func f() {
	{
		// defer 函数一定是在函数退出时才会执行，而不是代码块退出时执行
		defer fmt.Println("defer done")
		fmt.Println("code block")
	}

	fmt.Println("calling f")
}

// NOTE: WithClose

func WithClose(closer io.Closer, fn func()) {
	defer func() {
		closer.Close()
		fmt.Printf("close %s\n", closer.(*os.File).Name())
	}()
	fn()
}

func UseWithClose() {
	file, err := os.Open("data/foo.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	WithClose(file, func() {
		var content []byte
		content, err = io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(content))
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("with done")
}

// NOTE: 数据库事务

type Animal struct {
	Name string
}

// CreateAnimals 一个 GORM 官方文档中的示例 https://gorm.io/zh_CN/docs/transactions.html#手动事务
// func CreateAnimals(db *gorm.DB) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	if err := tx.Error; err != nil {
// 		return err
// 	}
//
// 	if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}
//
// 	if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}
//
// 	return tx.Commit().Error
// }

// CreateAnimals 优化后的写法
func CreateAnimals(db *gorm.DB) error {
	tx := db.Begin()
	defer tx.Rollback()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
		return err
	}

	if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// NOTE: 度量函数执行时间

func measureExecTime() func() {
	fmt.Println("calling measureExecTime")
	start := time.Now()
	return func() {
		fmt.Printf("execution use time: %s\n", time.Since(start))
	}
}

func fn() {
	defer measureExecTime()() // measureExecTime() 调用会同步执行

	fmt.Println("start")
	time.Sleep(2 * time.Second) // 模拟耗时操作
	fmt.Println("done")
}

// NOTE: defer 遇到 os.Exit 时不会被执行

// func f() {
// 	defer fmt.Println("deferred in f")
// 	fmt.Println("calling f")
// 	os.Exit(0)
// }

// NOTE: 一个过时的面试题

// func f() {
// 	for i := 0; i < 3; i++ {
// 		defer func() {
// 			fmt.Println(i)
// 		}()
// 	}
// }

// func f() {
// 	for i := 0; i < 3; i++ {
// 		defer fmt.Println(i)
// 	}
// }

// func f() {
// 	for i := 0; i < 3; i++ {
// 		defer func(i int) {
// 			fmt.Println(i)
// 		}(i)
// 	}
// }

// NOTE: defer nil 会引发 panic

func deferNil() {
	var f func()
	defer f()
	fmt.Println("calling deferNil")
}

func main() {
	// f()

	// fmt.Println(f())
	// fmt.Println(*f())

	// src := "data/file1.txt"
	// dst := "data/file2.txt"
	// written, err := CopyFile(dst, src)
	// fmt.Println(written, err)

	// processFile()
	// processFile1()
	// processFile2()

	UseWithClose()

	// printUser()

	// fn()

	// deferNil()
}
