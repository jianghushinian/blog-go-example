package calculator

// Add 返回两个整数的和（a + b）
func Add(a, b int) int {
	return a + b
}

// Subtract 返回两个整数的差（a - b）
func Subtract(a, b int) int {
	return a - b
}

// Multiply 返回两个整数的乘积（a * b）
func Multiply(a, b int) int {
	return a * b
}

// Divide 返回两个整数的商（a / b）
func Divide(a, b int) int {
	if b == 0 {
		return 0 // 简单处理除零错误
	}
	return a / b
}
