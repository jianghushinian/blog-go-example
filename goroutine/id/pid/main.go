package main

import (
	"fmt"
	"os"
)

func main() {
	// 获取当前进程的 id
	pid := os.Getpid()
	fmt.Println("process id:", pid)
}
