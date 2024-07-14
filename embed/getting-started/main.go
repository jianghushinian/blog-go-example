package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
)

// 适合嵌入单个文件（如配置数据、模板文件或一段文本）
//
//go:embed hello.txt
var content string

// 适合嵌入单个文件（如二进制文件：图片、字体或其他非文本数据）
//
//go:embed hello.txt
var contentBytes []byte

// 适合嵌入多个文件或整个目录（embed.FS 是只读文件系统）
//
//go:embed file
var fileFS embed.FS

//go:embed file/hello1.txt
//go:embed file/hello2.txt
var helloFS embed.FS

// 同上等价
// //go:embed file/hello1.txt file/hello2.txt
// var helloFS embed.FS

func main() {
	fmt.Printf("hello.txt content: %s\n", content)

	fmt.Printf("hello.txt content: %s\n", contentBytes)

	// NOTE: embed.FS 提供了 ReadFile 功能，可以直接读取文件内容，文件路径需要指明父目录 `file`
	hello1Bytes, _ := fileFS.ReadFile("file/hello1.txt")
	fmt.Printf("file/hello1.txt content: %s\n", hello1Bytes)

	// NOTE: embed.FS 提供了 ReadDir 功能，通过它可以遍历一个目录下的所有信息
	dir, _ := fs.ReadDir(fileFS, "file")
	for _, entry := range dir {
		info, _ := entry.Info()
		fmt.Printf("%+v\n", struct {
			Name  string
			IsDir bool
			Info  struct {
				Name string
				Size int64
				Mode fs.FileMode
			}
		}{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Info: struct {
				Name string
				Size int64
				Mode fs.FileMode
			}{Name: info.Name(), Size: info.Size(), Mode: info.Mode()},
		})
	}

	// NOTE: embed.FS 实现了 io/fs.FS 接口，可以返回它的子文件夹作为新的 io/fs.FS 文件系统
	subFS, _ := fs.Sub(helloFS, "file")
	hello2F, _ := subFS.Open("hello2.txt")
	hello2Bytes, _ := io.ReadAll(hello2F)
	fmt.Printf("file/hello2.txt content: %s\n", hello2Bytes)
}
