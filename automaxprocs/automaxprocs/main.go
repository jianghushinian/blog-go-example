package main

import (
	"fmt"
	"runtime"

	_ "go.uber.org/automaxprocs"
)

func main() {
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))
}
