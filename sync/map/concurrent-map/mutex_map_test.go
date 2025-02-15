package main

import (
	"fmt"
	"sync"
	"testing"
)

// NOTE: 这个测试不会终止
func TestRWMutexMap(t *testing.T) {
	m := NewRWMutexMap()

	// 并发读写 map
	go func() {
		for {
			m.Set("k", 1)
			fmt.Println("set k:", 1)
		}
	}()

	go func() {
		for {
			v, _ := m.Get("k")
			fmt.Println("read k:", v)
		}
	}()

	select {}
}

// NOTE: 这个测试不会终止
func TestRWMutexMapKV(t *testing.T) {
	m := NewRWMutexMapKV[string, int]()

	// 并发读写 map
	go func() {
		for {
			m.Set("k", 1)
			fmt.Println("set k:", 1)
		}
	}()

	go func() {
		for {
			v, _ := m.Get("k")
			fmt.Println("read k:", v)
		}
	}()

	select {}
}

// 测试读写锁 map 性能
func BenchmarkRWMutexMap(b *testing.B) {
	m := NewRWMutexMap()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Set(key, k)
			m.Get(key)
		}(i % 100000) // 使用 100000 个不同的 key
	}
	wg.Wait()
}

// 测试互斥锁 map 性能
func BenchmarkMutexMap(b *testing.B) {
	m := NewMutexMap()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Set(key, k)
			m.Get(key)
		}(i % 100000) // 使用 100000 个不同的 key
	}
	wg.Wait()
}

// 测试读写锁 map 性能（读多写少）
func BenchmarkRWMutexMapReadHeavy(b *testing.B) {
	m := NewRWMutexMap()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Get(string(rune(k)))
			}(j)
		}
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Set(string(rune(k)), k)
			}(j)
		}
	}
	wg.Wait()
}

// 测试互斥锁 map 性能（读多写少）
func BenchmarkMutexMapReadHeavy(b *testing.B) {
	m := NewMutexMap()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Get(string(rune(k)))
			}(j)
		}
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Set(string(rune(k)), k)
			}(j)
		}
	}
	wg.Wait()
}
