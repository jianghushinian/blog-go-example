package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 定义起始目录
	root := "./data"

	// 调用 Walk 函数遍历目录
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果发生错误，则输出错误并继续遍历
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		// 跳过名为 `.git` 的目录
		if info.IsDir() && info.Name() == ".git" {
			fmt.Printf("Skipping directory: %s\n", path)
			return filepath.SkipDir
		}

		// 跳过 Go 测试文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), "test.go") {
			fmt.Println("Skipping file:", path)
			return nil
		}

		// 输出每个文件或目录的路径
		fmt.Println(path)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v\n", err)
	}
}
