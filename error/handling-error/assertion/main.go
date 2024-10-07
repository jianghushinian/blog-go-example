package main

import (
	"fmt"
	"os"
)

func main() {
	// Type Assertion
	{
		// 尝试打开一个不存在的文件
		_, err := os.Open("nonexistent.txt")
		if err != nil {
			// 使用类型断言检查是否为 *os.PathError 类型
			if pathErr, ok := err.(*os.PathError); ok {
				fmt.Printf("Failed to %s file: %s\n", pathErr.Op, pathErr.Path)
				fmt.Println("Error message:", pathErr.Err)
			} else {
				// 其他类型的错误处理
				fmt.Println("Error:", err)
			}
		}
	}

	// Type Switch
	{
		// 尝试打开一个不存在的文件
		_, err := os.Open("nonexistent.txt")
		if err != nil {
			// 使用 switch type 检查错误类型
			switch e := err.(type) {
			case *os.PathError:
				fmt.Printf("Failed to %s file: %s\n", e.Op, e.Path)
				fmt.Println("Error message:", e.Err)
			default:
				// 其他类型的错误处理
				fmt.Println("Error:", err)
			}
		}
	}

	{
		// 尝试打开一个不存在的文件
		_, err := os.Open("nonexistent.txt")
		if err != nil {
			switch err.(type) {
			case *os.PathError, *os.LinkError:
				// do something
			default:
				// do something
			}
		}
	}
}
