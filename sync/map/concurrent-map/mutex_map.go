package main

import (
	"sync"
)

type MutexMap struct {
	mu sync.Mutex
	m  map[string]int
}

func NewMutexMap() *MutexMap {
	return &MutexMap{
		m: make(map[string]int),
	}
}

func (m *MutexMap) Set(key string, v int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = v
}

func (m *MutexMap) Get(key string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *MutexMap) Del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

func (m *MutexMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.m)
}

type RWMutexMap struct {
	rw sync.RWMutex
	m  map[string]int
}

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{
		m: make(map[string]int),
	}
}

func (m *RWMutexMap) Set(key string, v int) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[key] = v
}

func (m *RWMutexMap) Get(key string) (int, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *RWMutexMap) Del(key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.m, key)
}

func (m *RWMutexMap) Len() int {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return len(m.m)
}

type RWMutexMapKV[K comparable, V any] struct {
	rw sync.RWMutex
	m  map[K]V
}

// 作业：实现泛型版本 RWMutexMap

func NewRWMutexMapKV[K comparable, V any]() *RWMutexMapKV[K, V] {
	return &RWMutexMapKV[K, V]{
		m: make(map[K]V),
	}
}

func (m *RWMutexMapKV[K, V]) Set(k K, v V) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[k] = v
}

func (m *RWMutexMapKV[K, V]) Get(k K) (V, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, ok := m.m[k]
	return v, ok
}

func (m *RWMutexMapKV[K, V]) Del(k K) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.m, k)
}

func (m *RWMutexMapKV[K, V]) Len() int {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return len(m.m)
}
