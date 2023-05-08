package cmd

import (
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	cfgFile          string
	Verbose          bool
	Source           string
	Region           string
	author           string
	MarkdownDocs     bool
	ReStructuredDocs bool
	ManPageDocs      bool
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("hugo PersistentPreRun", args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("hugo PersistentPreRunE", args)
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("hugo PreRun", args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("hugo PreRunE", args)
		return nil
		// return errors.New("PreRunE err")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// NOTE: 生成文档时不执行命令
		if MarkdownDocs || ReStructuredDocs || ManPageDocs {
			return
		}
		fmt.Println("run hugo...")
		fmt.Printf("Verbose: %v\n", Verbose)
		fmt.Printf("Source: %v\n", Source)
		fmt.Printf("Region: %v\n", Region)
		fmt.Printf("Author: %v\n", viper.Get("author"))
		fmt.Printf("Config: %v\n", viper.AllSettings())
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("hugo PostRun", args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("hugo PersistentPostRun", args)
	},
	// TraverseChildren: true,
	// DisableSuggestions: true,
	// SuggestionsMinimumDistance: 1,
	Example: "hugo example", // 使用示例，将在 `hugo -h` 时显示
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GenDocs() {
	var err error
	switch {
	case MarkdownDocs:
		err = doc.GenMarkdownTree(rootCmd, "./docs/md")
	case ReStructuredDocs:
		err = doc.GenReSTTree(rootCmd, "./docs/rest")
	case ManPageDocs:
		t := time.Now()
		header := &doc.GenManHeader{
			Title:   "hugo",
			Section: "1",
			Manual:  "hugo Manual",
			Source:  "hugo source",
			Date:    &t,
		}
		err = doc.GenManTree(rootCmd, header, "./docs/man")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 初始化配置
	cobra.OnInitialize(initConfig)
	// 使用本地标志指定配置
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")

	// 持久标志
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	// 本地标志
	rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")

	// 必选标志
	rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
	if err := rootCmd.MarkFlagRequired("region"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 绑定标志到 Viper
	rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
	if err := viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 生成文档
	rootCmd.Flags().BoolVarP(&MarkdownDocs, "md-docs", "m", false, "gen Markdown docs")
	rootCmd.Flags().BoolVarP(&ReStructuredDocs, "rest-docs", "t", false, "gen ReStructured Text docs")
	rootCmd.Flags().BoolVarP(&ManPageDocs, "man-docs", "a", false, "gen Man Page docs")

	// 定制使用 `help` 命令查看帮助信息输出结果
	/*
		rootCmd.SetHelpCommand(&cobra.Command{
			Use:    "help",
			Short:  "Custom help command",
			Hidden: true,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Custom help command")
			},
		})
	*/
	// 定制使用 `-h/--help` 命令查看帮助信息输出结果
	/*
		rootCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
			fmt.Println(strings)
		})
	*/
	// 定制帮助信息模板
	/*
			rootCmd.SetHelpTemplate(`Custom Help Template:
		Usage:
			{{.UseLine}}
		Description:
			{{.Short}}
		Commands:
		{{- range .Commands}}
			{{.Name}}: {{.Short}}
		{{- end}}
		`)
	*/

	// 定制 Usage Message
	/*
		rootCmd.SetUsageFunc(func(command *cobra.Command) error {
			fmt.Printf("Custom usage for command: %s\n", command.Name())
			return nil
		})
	*/
	// Usage 模板
	/*
			rootCmd.SetUsageTemplate(`Custom Usage Template:
		Usage: {{.CommandPath}} [command]

		Description: {{.Short}}

		Available Commands:
		{{- range .Commands -}}
			{{- if not .Hidden -}}
				{{rpad .Name .NamePadding}}{{.Short}}{{ "\n" }}
			{{- end -}}
		{{- end }}
		`)
	*/
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
