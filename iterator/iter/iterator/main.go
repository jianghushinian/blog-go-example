// NOTE: Go 1.23 迭代器

package main

import "fmt"

// NOTE: 最简单的迭代器

/*
func iterator(yield func() bool) {
	for i := 0; i < 5; i++ {
		if !yield() {
			return
		}
	}
}

func main() {
	i := 0
	for range iterator {
		fmt.Printf("i=%d\n", i)
		i++
	}
}
*/

// NOTE: 控制迭代次数

/*
func iterator(n int) func(yield func() bool) {
	return func(yield func() bool) {
		for i := 0; i < n; i++ {
			if !yield() {
				return
			}
		}
	}
}

func main() {
	i := 0
	for range iterator(3) {
		fmt.Printf("i=%d\n", i)
		i++
	}
}
*/

// NOTE: 输出一个值

/*
func iterator(n int) func(yield func(v int) bool) {
	return func(yield func(v int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	i := 0
	for v := range iterator(10) {
		if i >= 5 {
			break
		}
		fmt.Printf("%d => %d\n", i, v)
		i++
	}
}
*/

// NOTE: 输出两个值

/*
func iterator(slice []int) func(yield func(i, v int) bool) {
	return func(yield func(i int, v int) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

func main() {
	s := []int{0, 1, 2, 3, 4}
	for i, v := range iterator(s) {
		if i == 2 {
			continue
		}
		fmt.Printf("%d => %d\n", i, v)
	}
}
*/

// NOTE: map 迭代器

func iterator(m map[string]int) func(yield func(k string, v int) bool) {
	return func(yield func(k string, v int) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func main() {
	m := map[string]int{
		"a": 0,
		"b": 1,
		"c": 2,
	}
	for k, v := range iterator(m) {
		fmt.Printf("%s: %d\n", k, v)
	}
}
