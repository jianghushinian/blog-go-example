package foo

import (
	"time"
	// 被拉取的包需要显式导入（除了 runtime 包）

	_ "fmt"
	_ "unsafe"
	_ "github.com/jianghushinian/blog-go-example/directive/linkname/bar"
)

// Pull 模式（拉取外部实现）

//go:linkname Add github.com/jianghushinian/blog-go-example/directive/linkname/bar.add
func Add(a, b int) int

func Div(a, b int) int

// Handshake 模式（双方握手模式）

//go:linkname Hello github.com/jianghushinian/blog-go-example/directive/linkname/bar.hello
func Hello(name string) string

// 标准库默认不允许链接
// 编译时不指定 -ldflags=-checklinkname=0 则会报错

//go:linkname FooPrintln fmt.Println
func FooPrintln(a ...any) (n int, err error)

//go:linkname TooLarge 	fmt.tooLarge
func TooLarge(x int) bool

// 但是可以直接引用使用了 //go:linkname 标记的符号

//go:linkname Now time.Now
func Now() time.Time

// linkname 不止可以链接函数，实际上可以链接任何函数或变量

//go:linkname X github.com/jianghushinian/blog-go-example/directive/linkname/bar.x
var X int
