package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	var (
		mu  sync.Mutex
		out bytes.Buffer
		wg  sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			var local bytes.Buffer
			for j := 0; j < 1000; j++ {
				local.WriteString("x")
			}

			mu.Lock()
			out.Write(local.Bytes())
			mu.Unlock()
		})
	}
	wg.Wait()
	fmt.Println(out.Len())
}
