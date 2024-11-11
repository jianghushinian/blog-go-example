package once

import "sync"

func main() {
	var once sync.Once
	once.Do(func() {
		println("Only once")
	})
}
