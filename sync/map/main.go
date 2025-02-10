package main

import (
	"fmt"
	"sync"
)

func main() {
	var s sync.Map

	// 存储键值对
	s.Store("name", "江湖十年")
	s.Store("age", 20)
	s.Store("location", "Beijing")

	// 读取值
	if value, ok := s.Load("name"); ok {
		fmt.Println("name:", value)
	}

	// 删除一个键
	s.Delete("age")

	// 遍历 sync.Map
	s.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %s\n", key, value)
		return true // 继续遍历
	})
}
