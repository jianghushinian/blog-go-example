package main

import (
	"fmt"
	"strconv"
	"time"
)

// MagicCalculator 魔术计算器
type MagicCalculator struct {
	targetTime int    // 目标时间转换的数字
	timestamp  string // 实际时间字符串
}

// NewMagicCalculator 创建一个魔术计算器实例
func NewMagicCalculator() *MagicCalculator {
	// 获取当前时间
	now := time.Now()

	// 生成类似 "2162227" 的时间数字
	// 格式: 月(1-2 位) + 日(2 位) + 小时(2 位) + 分钟(2 位)
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()

	// 构建时间字符串和数字
	timestamp := fmt.Sprintf("%d%02d%02d%02d", month, day, hour, minute)

	// 转换为整数
	target, _ := strconv.ParseInt(timestamp, 10, 64)

	return &MagicCalculator{
		targetTime: int(target),
		timestamp:  timestamp,
	}
}

// GetMagicNumber 计算魔术数字
func (mc *MagicCalculator) GetMagicNumber(num1, num2 int) int {
	// 魔术公式: target - (num1 + num2)
	return mc.targetTime - (num1 + num2)
}

func InteractiveMagic() {
	fmt.Println("=== 交互式魔术体验 ===")
	fmt.Println("请按照提示输入数字，我会展示魔术的原理")

	mc := NewMagicCalculator()

	var num1, num2 int
	fmt.Print("请输入第一个数: ")
	fmt.Scan(&num1)
	fmt.Print("请输入第二个数: ")
	fmt.Scan(&num2)

	fmt.Printf("\n你输入的是: %d 和 %d\n", num1, num2)

	magicNum := mc.GetMagicNumber(num1, num2)
	fmt.Printf("魔术数字（第三个数）是: %d\n", magicNum)

	fmt.Printf("\n验证: %d + %d + %d = %d\n", num1, num2, magicNum, mc.targetTime)
	fmt.Printf("这个数字代表的时间是: %s\n", mc.timestamp)
}

func main() {
	InteractiveMagic()
}
