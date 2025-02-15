package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: 这个测试不会终止
func TestShardingMap(t *testing.T) {
	m := NewShardingMap(10)

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

// 并发读写正确性
func TestConcurrentSetGet(t *testing.T) {
	m := NewShardingMap(16)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("key%d", k), k)
			val, ok := m.Get(fmt.Sprintf("key%d", k))
			assert.True(t, ok)
			assert.Equal(t, k, val)
		}(i)
	}
	wg.Wait()
}

// 删除操作验证
func TestDelete(t *testing.T) {
	m := NewShardingMap(16)
	m.Set("key1", 42)
	m.Del("key1")
	val, ok := m.Get("key1")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// 哈希冲突处理
func TestHashCollision(t *testing.T) {
	m := NewShardingMap(1) // 强制所有 key 到同一分片
	m.Set("key1", 100)
	m.Set("key2", 200)
	val1, _ := m.Get("key1")
	val2, _ := m.Get("key2")
	assert.Equal(t, 100, val1)
	assert.Equal(t, 200, val2)
}

// 测试分片 map 性能
func BenchmarkShardingMap(b *testing.B) {
	m := NewShardingMap(32)
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

// 对比 sync.Map
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Store(key, k)
			m.Load(key)
		}(i % 100000) // 使用 100000 个不同的 key
	}
	wg.Wait()
}

// 测试分片 map 性能（读多写少）
func BenchmarkShardingMapReadHeavy(b *testing.B) {
	m := NewShardingMap(32)
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

// 对比 sync.Map（读多写少）
func BenchmarkSyncMapReadHeavy(b *testing.B) {
	var m sync.Map
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Load(string(rune(k)))
			}(j)
		}
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				m.Store(string(rune(k)), k)
			}(j)
		}
	}
	wg.Wait()
}
