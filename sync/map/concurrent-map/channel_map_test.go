package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: 这个测试不会终止
func TestChannelMap(t *testing.T) {
	m := NewChannelMap()

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

// 测试 channel map 性能
func BenchmarkChannelMap(b *testing.B) {
	m := NewChannelMap()
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

// 测试 channel map 性能（读多写少）
func BenchmarkChannelMapReadHeavy(b *testing.B) {
	m := NewChannelMap()
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

func TestChannelMapKV(t *testing.T) {
	m := NewChannelMapKV[string, string]()

	m.Set("k", "v")
	v, ok := m.Get("k")
	assert.Equal(t, true, ok)
	assert.Equal(t, "v", v)

	l := m.Len()
	assert.Equal(t, 1, l)

	m.Del("k")
	v, ok = m.Get("k")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", v)

	l = m.Len()
	assert.Equal(t, 0, l)
}
