// Package calculator provides math operations for financial calculations.
//
// 本包实现安全的四则运算，包含加减乘除功能。
// 所有函数均针对整型设计，适用于基础计算场景。
//
// 使用示例：
//
//	sum := calculator.Add(5, 3)          // 8
//	diff := calculator.Subtract(5, 3)    // 2
//	product := calculator.Multiply(5, 3) // 15
//	quotient := calculator.Divide(6, 3)  // 2
//
// 注意事项：
//  1. 除法函数 Divide() 在除数为 0 时返回 0
//  2. 整数除法会截断小数部分（如 5/2=2）
package calculator // import "govanityurls.jianghushinian.cn/blog-go-example/doc/docgo/calculator"
