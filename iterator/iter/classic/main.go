// NOTE: 手写迭代器

package main

// NOTE: 迭代器模式

type Iterator struct {
	data  []int
	index int
}

func NewIterator(data []int) *Iterator {
	return &Iterator{data: data, index: 0}
}

func (it *Iterator) HasNext() bool {
	return it.index < len(it.data)
}

func (it *Iterator) Next() int {
	if !it.HasNext() {
		panic("Stop iteration")
	}
	value := it.data[it.index]
	it.index++
	return value
}

// func main() {
// 	it := NewIterator([]int{0, 1, 2, 3, 4})
// 	for it.HasNext() {
// 		fmt.Println(it.Next())
// 	}
// }

// NOTE: 回调函数风格迭代器

// func main() {
// 	// 循环链表
// 	r := ring.New(5)
// 	// 初始化链表
// 	for i := 0; i < r.Len(); i++ {
// 		r.Value = i  // 为当前节点赋值
// 		r = r.Next() // 移动到下一个节点
// 	}
// 	// 迭代器
// 	r.Do(func(v any) {
// 		fmt.Println(v)
// 	})
// }

// NOTE: Go 风格迭代器

func generator(n int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

// func main() {
// 	for n := range generator(5) {
// 		fmt.Println(n)
// 	}
// }
