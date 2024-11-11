package once

import "github.com/jianghushinian/blog-go-example/sync/once/inlining/myonce/sync"

func main() {
	var once sync.Once
	once.Do(func() {
		println("Only once")
	})
}
