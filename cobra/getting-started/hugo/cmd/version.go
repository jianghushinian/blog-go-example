package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("version PersistentPreRun")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("version PreRun")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("version PostRun")
	// },
	// PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("version PersistentPostRun")
	// },
	Args: cobra.MaximumNArgs(2), // 使用内置的验证函数，位置参数多于 2 个则报错
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// 自定义帮助命令
	/*
		versionCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
			fmt.Println("version help")
		})
	*/
	// 定制帮助信息模板
	/*
			versionCmd.SetHelpTemplate(`Custom Help Template:
		Usage:
			{{.UseLine}}
		Description:
			{{.Short}}
		`)
	*/
}
