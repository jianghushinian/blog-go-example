package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lscolor",
	Short: "lscolor lists files and directories with colors",
	Args:  cobra.MaximumNArgs(1), // 允许零个或一个参数
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	dir := "." // 默认使用当前目录
	if len(args) > 0 {
		dir = args[0] // 如果提供了参数，则使用该参数
	}

	// 获取当前目录的文件和子目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		color.Red("error reading directory: %s", err)
		return
	}

	// 遍历并输出每个条目
	for _, entry := range entries {
		if entry.IsDir() {
			color.Cyan(entry.Name()) // 目录用蓝绿色
		} else {
			color.White(entry.Name()) // 文件用白色
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
