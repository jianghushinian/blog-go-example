// NOTE: 泛型版本

package main

// NOTE: 输出零个值

type Seq0 func(yield func() bool)

func iter0[Slice ~[]E, E any](s Slice) Seq0 {
	return func(yield func() bool) {
		for range s {
			if !yield() {
				return
			}
		}
	}
}

/*
func main() {
	s1 := []int{1, 2, 3}
	i := 0
	for range iter0(s1) {
		fmt.Printf("i=%d\n", i)
		i++
	}

	fmt.Println("--------------")

	s2 := []string{"a", "b", "c"}
	i = 0
	for range iter0(s2) {
		fmt.Printf("i=%d\n", i)
		i++
	}
}
*/

// NOTE: 输出一个值

type Seq1[V any] func(yield func(V) bool)

func iter1[Slice ~[]E, E any](s Slice) Seq1[E] {
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

/*
func main() {
	s1 := []int{1, 2, 3}
	for v := range iter1(s1) {
		fmt.Printf("v=%d\n", v)
	}

	fmt.Println("--------------")

	s2 := []string{"a", "b", "c"}
	for v := range iter1(s2) {
		fmt.Printf("v=%s\n", v)
	}
}
*/

// NOTE: 输出两个值

type Seq2[K, V any] func(yield func(K, V) bool)

func iter2[Slice ~[]E, E any](s Slice) Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}

/*
func main() {
	s1 := []int{1, 2, 3}
	for i, v := range iter2(s1) {
		fmt.Printf("%d=%d\n", i, v)
	}

	fmt.Println("--------------")

	s2 := []string{"a", "b", "c"}
	for i, v := range iter2(s2) {
		fmt.Printf("%d=%s\n", i, v)
	}
}
*/
