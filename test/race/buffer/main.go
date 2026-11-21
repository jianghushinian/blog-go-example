package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	var (
		buf bytes.Buffer
		wg  sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				buf.WriteString("x")
			}
		})
	}
	wg.Wait()
	fmt.Println(buf.Len())
}
