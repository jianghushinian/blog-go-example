package main

import (
	"hash/maphash"
	"sync"
)

var seed = maphash.MakeSeed()

func hashKey(key string) uint64 {
	return maphash.String(seed, key)
}

type ShardingMap struct {
	locks  []sync.RWMutex
	shards []map[string]int
}

func NewShardingMap(size int) *ShardingMap {
	sm := &ShardingMap{
		locks:  make([]sync.RWMutex, size),
		shards: make([]map[string]int, size),
	}
	for i := 0; i < size; i++ {
		sm.shards[i] = make(map[string]int)
	}
	return sm
}

func (m *ShardingMap) getShardIdx(key string) uint64 {
	hash := hashKey(key)
	return hash % uint64(len(m.shards))
}

func (m *ShardingMap) Set(key string, value int) {
	idx := m.getShardIdx(key)
	m.locks[idx].Lock()
	defer m.locks[idx].Unlock()
	m.shards[idx][key] = value
}

func (m *ShardingMap) Get(key string) (int, bool) {
	idx := m.getShardIdx(key)
	m.locks[idx].RLock()
	defer m.locks[idx].RUnlock()
	value, ok := m.shards[idx][key]
	return value, ok
}

func (m *ShardingMap) Del(key string) {
	idx := m.getShardIdx(key)
	m.locks[idx].Lock()
	defer m.locks[idx].Unlock()
	delete(m.shards[idx], key)
}

// Len 遍历所有分片并累加长度，但在高并发场景下，返回的是瞬时近似值，而非精确值
func (m *ShardingMap) Len() int {
	l := 0
	// 遍历过程中，可能会有人修改遍历过的 shard，所以长度可能不精确
	for i, shard := range m.shards {
		m.locks[i].RLock()
		l += len(shard)
		m.locks[i].RUnlock()
	}
	return l
}

// NOTE: 这只是一个分片 map 的 demo 实现，可以优化的点包括但不限于
// 1. 可以在 ShardingMap 结构体中维护 seed，每个 map 实例都使用自己的 seed
// 2. 可以在 ShardingMap 结构体中维护一个 len 字段，用来记录 map 的总长度，这样每次操作时需要记得同步修改 len 字段的值
// 3. 可以实现更多常用方法
// 4. 可以选用一个性能更好的 hash 计算函数
// 5. 支持泛型，对 hash 函数要求更高，需要一个可以为 any 类型计算哈希值的函数
