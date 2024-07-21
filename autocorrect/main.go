package main

import (
	"fmt"

	"github.com/longbridgeapp/autocorrect"
)

func main() {
	fmt.Println(autocorrect.Format("长桥LongBridge App下载"))
	// => "长桥 LongBridge App 下载"

	fmt.Println(autocorrect.Format("Ruby 2.7版本第1次发布"))
	// => "Ruby 2.7 版本第 1 次发布"

	fmt.Println(autocorrect.Format("于3月10日开始"))
	// => "于 3 月 10 日开始"

	fmt.Println(autocorrect.Format("包装日期为2013年3月10日"))
	// => "包装日期为 2013 年 3 月 10 日"

	fmt.Println(autocorrect.Format("生产环境中使用Go"))
	// => "生产环境中使用 Go"

	fmt.Println(autocorrect.Format("本番環境でGoを使用する"))
	// => "本番環境で Go を使用する"

	fmt.Println(autocorrect.Format("프로덕션환경에서Go사용"))
	// => "프로덕션환경에서 Go 사용"

	fmt.Println(autocorrect.Format("需要符号?自动转换全角字符、数字:我们将在１６：３２分出发去ＣＢＤ中心."))
	// => "需要符号？自动转换全角字符、数字：我们将在 16:32 分出发去 CBD 中心。"
}
