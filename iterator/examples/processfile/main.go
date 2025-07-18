package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"strings"
)

// NOTE: 以一个读取文件内容的函数为例

// 实现一：一次性加载整个文件，可能出现内存溢出

func ProcessFile1(filename string) {
	data, _ := os.ReadFile(filename)
	lines := strings.Split(string(data), "\n") // 按换行切分
	for i, line := range lines {
		fmt.Printf("line %d: %s\n", i, line)
	}
}

// 实现二：使用 bufio 迭代器实现

func ProcessFile2(filename string) {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		fmt.Printf("line %d: %s\n", i, scanner.Text())
		i++
	}
}

// 实现三：使用 Go 1.23 迭代器进行重构

// NOTE: 实现者

func ReadLines(filename string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		file, _ := os.Open(filename)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			if !yield(i, scanner.Text()) { // 按需生成
				return
			}
			i++
		}
	}
}

// NOTE: 使用者
// 把复杂留给实现着，把**标准**留个使用者

func ProcessFile3(filename string) {
	for i, line := range ReadLines(filename) {
		fmt.Printf("line %d: %s\n", i, line)
	}
}

func main() {
	filename := "demo/main.go"
	ProcessFile1(filename)
	fmt.Println("--------------")
	ProcessFile2(filename)
	fmt.Println("--------------")
	ProcessFile3(filename)
}
