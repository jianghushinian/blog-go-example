package main

import "fmt"

// 生成所有排列组合
func permute(items []string, start int) {
	if start == len(items) {
		// 当排列完成时，输出当前的排列
		// fmt.Println(items)
		magic(items)
		return
	}

	for i := start; i < len(items); i++ {
		// 交换位置
		items[start], items[i] = items[i], items[start]
		// 递归调用
		permute(items, start+1)
		// 交换回原来的位置，回溯
		items[start], items[i] = items[i], items[start]
	}
}

// 魔术
func magic(items []string) {
	old := make([]string, len(items))
	copy(old, items)

	// 1. 筷子跟它左边的物品互换，如果筷子已经在最左边，则无需移动
	for i := 1; i < len(items); i++ {
		if items[i] == "筷子🥢" {
			// 筷子如果不在最左边，交换到最左边
			items[i], items[0] = items[0], items[i]
			break
		}
	}

	// 2. 杯子跟它右边的物品互换，如果杯子已经在最右边，则无需移动
	for i := len(items) - 2; i >= 0; i-- {
		if items[i] == "杯子🍺" {
			// 杯子如果不在最右边，交换到最右边
			items[i], items[len(items)-1] = items[len(items)-1], items[i]
			break
		}
	}

	// 3. 勺子跟它左边的物品互换，如果勺子已经在最左边，则无需移动
	for i := 1; i < len(items); i++ {
		if items[i] == "勺子🥄" {
			// 勺子如果不在最左边，交换到最左边
			items[i], items[0] = items[0], items[i]
			break
		}
	}

	// 打印当前和经过魔术操作后的排列
	fmt.Println("当前排列：", old, " => ", "魔术操作后：", items)
}

func main() {
	items := []string{"筷子🥢", "杯子🍺", "勺子🥄"}
	permute(items, 0)
}
