package abs

import "fmt"

// --------------示例测试---------------

func ExampleAbs() {
	fmt.Println(Abs(-1))
	fmt.Println(Abs(2))
	// Output:
	// 1
	// 2
}

func ExampleAbs_unordered() {
	fmt.Println(Abs(2))
	fmt.Println(Abs(-1))
	// Unordered Output:
	// 1
	// 2
}
