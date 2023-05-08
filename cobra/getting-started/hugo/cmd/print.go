package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var printFlag string

var printCmd = &cobra.Command{
	Use: "print [OPTIONS] [COMMANDS]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run print...")
		fmt.Printf("printFlag: %v\n", printFlag)
		fmt.Printf("Source: %v\n", Source)
		// 命令行位置参数列表：例如执行 `hugo print a b c d` 将得到 [a b c d]
		fmt.Printf("args: %v\n", args)
	},
	// 使用自定义参数验证函数
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		if len(args) > 4 {
			return errors.New("the number of args cannot exceed 4")
		}
		if args[0] != "a" {
			return errors.New("first argument must be 'a'")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	// 本地标志
	printCmd.Flags().StringVarP(&printFlag, "flag", "f", "", "print flag for local")
}
