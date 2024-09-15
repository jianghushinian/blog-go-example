package main

import (
	"fmt"
	"runtime"
	"strconv"
)

func printStack(skip int) {
	var pcs [30]uintptr
	n := runtime.Callers(skip, pcs[:])

	for i := 0; i < n; i++ {
		pc := pcs[i]
		fn := runtime.FuncForPC(pc - 1)
		file, line := fn.FileLine(pc - 1)
		fmt.Printf("Func Name: %s\n", fn.Name())
		fmt.Printf("File: %s, Line: %s\n\n", file, strconv.Itoa(line))
	}
}

func Print(skip int) {
	printStack(skip)
}

func main() {
	Print(0)

	fmt.Println("============================================")

	Print(3)
}

func printStackByCallersFrames(skip int) {
	var pcs [30]uintptr
	n := runtime.Callers(skip, pcs[:])

	fs := runtime.CallersFrames(pcs[:n])
	for {
		f, ok := fs.Next()
		fmt.Printf("Func Name: %s\n", f.Function)
		fmt.Printf("File: %s, Line: %s\n\n", f.File, strconv.Itoa(f.Line))
		if !ok {
			break
		}
	}
}

/*
$ go run main.go
Func Name: runtime.Callers
File: /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/extern.go, Line: 325

Func Name: main.printStack
File: /go/blog-go-example/error/pkg-errors/runtime-example/main.go, Line: 11

Func Name: main.Print
File: /go/blog-go-example/error/pkg-errors/runtime-example/main.go, Line: 23

Func Name: main.main
File: /go/blog-go-example/error/pkg-errors/runtime-example/main.go, Line: 27

Func Name: runtime.main
File: /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/proc.go, Line: 271

Func Name: runtime.goexit
File: /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/asm_arm64.s, Line: 1222

============================================
Func Name: main.main
File: /go/blog-go-example/error/pkg-errors/runtime-example/main.go, Line: 31

Func Name: runtime.main
File: /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/proc.go, Line: 271

Func Name: runtime.goexit
File: /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.7.darwin-arm64/src/runtime/asm_arm64.s, Line: 1222
*/
