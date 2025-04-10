package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))
}
