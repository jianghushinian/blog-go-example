package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 我们使用 `int` 类型的切片 `pokers` 来存储刘谦手中的 4 张扑克牌
	var pokers []int = []int{2, 7, 6, 5}

	// 1. 洗牌，随意打乱
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(pokers), func(i, j int) {
		pokers[i], pokers[j] = pokers[j], pokers[i]
	})
	fmt.Printf("1. 洗牌后的牌：%v\n", pokers)

	// 2. 对折，然后撕开，即切片内容 * 2
	pokers = append(pokers, pokers...)
	fmt.Printf("2. 对折后的牌：%v\n", pokers)

	// 3. 问问自己名字有几个字，就从最上面拿出对应个数的牌放到底部，例如「刘谦」名字有 2 个字，即将切片前 2 个元素取出，放到切片最后
	pokers = append(pokers[2:], pokers[:2]...)
	fmt.Printf("3. 问问自己名字有几个字，就从最上面拿出对应个数的牌放到底部：%v\n", pokers)

	// 4. 拿起最上面的 3 张牌，插入中间任意位置，这里将切片前 3 个元素取出，并将其插入最后一个元素之前
	pokers = append(pokers[3:7], pokers[0], pokers[1], pokers[2], pokers[7])
	fmt.Printf("4. 拿起最上面的 3 张牌，插入中间任意位置: %v\n", pokers)

	// 5. 拿出最上面的 1 张牌，藏于秘密的地方，比如屁股下，这里使用 top 变量暂存
	top := pokers[0]
	pokers = pokers[1:]
	fmt.Printf("5. 拿出最上面的 1 张牌：%d, %v\n", top, pokers)

	// 6. 如果你是南方人，从上面拿起 1 张牌；如果你是北方人，则从上面拿起 2 张牌；假如我们不确定自己是南方人还是北方人，那就干脆拿起 3 张牌，然后插入中间任意位置
	pokers = append(pokers[2:6], pokers[0], pokers[1], pokers[6])
	fmt.Printf("6. 从上面拿出 2 张牌: %v\n", pokers)

	// 7. 如果你是男生，从上面拿起 1 张牌；如果你是女生，则从上面拿起 2 张牌，撒到空中（扔掉）
	pokers = pokers[1:]
	fmt.Printf("7. 如果你是男生，从上面拿起 1 张牌；如果你是女生，则从上面拿起 2 张牌，撒到空中（扔掉）：%v\n", pokers)

	// 8. 魔法时刻，在遥远的魔术的历史上，流传了一个七字真言「见证奇迹的时刻」，可以带给我们幸福。现在，我们每念一个字，从上面拿一张放到最底部，即需要完成 7 次同样的操作
	//    我们可以用一个 `for loop` 实现
	for range []string{"见", "证", "奇", "迹", "的", "时", "刻"} {
		pokers = append(pokers[1:], pokers[0])
	}
	fmt.Printf("8. 见证奇迹的时刻：%v\n", pokers)

	// 9. 最后一个环节，叫「好运留下来，烦恼丢出去」，在念到「好运留下来」时，从上面拿起 1 张牌放入底部；在念到「烦恼丢出去」时，从上面拿起 1 张牌扔掉，女生需要完成 4 次同样的操作，男生需要完成 5 次同样的操作
	//    同样可以用一个 `for loop` 实现
	for range []int{1, 2, 3, 4, 5} {
		// 好运留下来
		pokers = append(pokers[1:], pokers[0])
		// 烦恼丢出去
		pokers = pokers[1:]
	}
	fmt.Printf("9. 好运留下来，烦恼丢出去：%v\n", pokers)

	// 最后，我们将见证奇迹：
	fmt.Printf("见证奇迹：%d == %d", top, pokers[0])
}
