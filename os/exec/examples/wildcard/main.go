package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// NOTE: 当命令参数中带有通配符

// func main() {
// 	cmd := exec.Command("ls", "-l", "/var/log/*.log")
//
// 	output, err := cmd.CombinedOutput() // 获取标准输出和错误输出
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output))
// }

// NOTE: 使用 bash -c 来解析通配符

// func main() {
// 	// 使用 bash -c 来解析通配符
// 	cmd := exec.Command("bash", "-c", "ls -l /var/log/*.log")
//
// 	output, err := cmd.CombinedOutput() // 获取标准输出和错误输出
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output))
// }

// NOTE: 使用 Go 标准库中的 filepath.Glob 来手动解析通配符

func main() {
	// 匹配通配符路径
	files, err := filepath.Glob("/var/log/*.log")
	if err != nil {
		log.Fatalf("Glob failed: %v", err)
	}
	if len(files) == 0 {
		log.Println("No matching files found")
		return
	}

	// 将匹配到的文件传给 ls 命令
	args := append([]string{"-l"}, files...)
	cmd := exec.Command("ls", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
