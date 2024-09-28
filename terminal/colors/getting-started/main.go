package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	// 标准颜色输出
	{
		color.Cyan("Prints text in cyan.")
		color.Blue("Prints %s in blue.", "text")
		color.Red("Prints text in red.")
		color.Magenta("And many others...")
	}

	// 混合使用多种属性
	{
		// 创建一个 color 对象，输出效果：蓝绿色 + 下划线
		c := color.New(color.FgCyan).Add(color.Underline)
		c.Println("Prints cyan text with an underline.") // 注意 Println 输出自动加换行

		// 输出效果：蓝绿色 + 加粗
		d := color.New(color.FgCyan, color.Bold)
		d.Printf("This prints bold cyan %s\n", "too!.") // 注意 Printf 需要手动加换行

		// 输出效果：红色 + 白色背景
		red := color.New(color.FgRed)
		whiteBackground := red.Add(color.BgWhite)
		whiteBackground.Println("Red text with white background.")
	}

	// 输出到指定位置
	{
		f, _ := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		color.New(color.FgBlue).Fprintln(f, "blue color!")

		blue := color.New(color.FgBlue)
		blue.Fprint(f, "This will print text in blue.\n")
	}

	// 混合输出普通字符和带颜色的字符
	{
		// 输出效果：普通文本 + 黄色 warning + 红色 error
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

		// 模拟输出日志内容
		fmt.Println(color.GreenString("Info:"), "a info log message")
		fmt.Printf("%v: %v\n", color.RedString("Warn"), "a warning log message")
	}

	// 改变原生代码
	{
		// 修改标准输出
		color.Set(color.FgYellow)
		fmt.Println("Existing text will now be in yellow")
		fmt.Printf("This one %s\n", "too")
		color.Unset() // 恢复设置
		fmt.Println("This is normal text")

		func() {
			// 在函数中使用
			color.Set(color.FgMagenta, color.Bold)
			defer color.Unset()

			fmt.Println("All text will now be bold magenta.")
		}()
		fmt.Println("This is normal text too")
	}
}
