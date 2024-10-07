package main

import (
	"errors"
	"fmt"
)

func main() {
	var userID int

	// 创建一个错误值
	err1 := errors.New("example err1")
	// 格式化错误消息
	err2 := fmt.Errorf("example err2: %d", userID)

	fmt.Println(err1, err2)
}
