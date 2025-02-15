package main

type ChannelMap struct {
	cmd chan command
	m   map[string]int
}

type command struct {
	action string // "get", "set", "delete"
	key    string
	value  int
	result chan<- result
}

type result struct {
	value int
	ok    bool
}

func NewChannelMap() *ChannelMap {
	sm := &ChannelMap{
		cmd: make(chan command),
		m:   make(map[string]int),
	}
	go sm.run()
	return sm
}

func (m *ChannelMap) run() {
	for cmd := range m.cmd {
		switch cmd.action {
		case "get":
			value, ok := m.m[cmd.key]
			cmd.result <- result{value, ok}
		case "set":
			m.m[cmd.key] = cmd.value
		case "delete":
			delete(m.m, cmd.key)
		}
	}
}

func (m *ChannelMap) Set(key string, value int) {
	m.cmd <- command{action: "set", key: key, value: value}
}

func (m *ChannelMap) Get(key string) (int, bool) {
	res := make(chan result)
	m.cmd <- command{action: "get", key: key, result: res}
	r := <-res
	return r.value, r.ok
}

func (m *ChannelMap) Del(key string) {
	m.cmd <- command{action: "delete", key: key}
}

type ChannelMapKV[K comparable, V any] struct {
	cmd  chan commandKV[K, V]
	data map[K]V
}

type commandKV[K comparable, V any] struct {
	action string
	key    K
	value  V
	result chan<- resultKV[K, V]
}

type resultKV[K comparable, V any] struct {
	value V
	ok    bool
	len   int
}

func NewChannelMapKV[K comparable, V any]() *ChannelMapKV[K, V] {
	cm := &ChannelMapKV[K, V]{
		cmd:  make(chan commandKV[K, V]),
		data: make(map[K]V),
	}
	go cm.run()
	return cm
}

func (m *ChannelMapKV[K, V]) run() {
	for cmd := range m.cmd {
		switch cmd.action {
		case "get":
			value, ok := m.data[cmd.key]
			cmd.result <- resultKV[K, V]{value, ok, 0}
		case "set":
			m.data[cmd.key] = cmd.value
		case "delete":
			delete(m.data, cmd.key)
		case "len":
			length := len(m.data)
			cmd.result <- resultKV[K, V]{*new(V), false, length}
		}
	}
}

func (m *ChannelMapKV[K, V]) Set(key K, value V) {
	m.cmd <- commandKV[K, V]{action: "set", key: key, value: value}
}

func (m *ChannelMapKV[K, V]) Get(key K) (V, bool) {
	res := make(chan resultKV[K, V])
	m.cmd <- commandKV[K, V]{action: "get", key: key, result: res}
	r := <-res
	return r.value, r.ok
}

func (m *ChannelMapKV[K, V]) Del(key K) {
	m.cmd <- commandKV[K, V]{action: "delete", key: key}
}

func (m *ChannelMapKV[K, V]) Len() int {
	res := make(chan resultKV[K, V])
	m.cmd <- commandKV[K, V]{action: "len", result: res}
	r := <-res
	return r.len
}

// NOTE: 对于 channel 的 map 实现当然是不建议的，任何技术都有其适用场景，在这里用不太合适
// 如果非要这样使用，记得给 ChannelMap 提供一个 Stop 方法来停止内部的 goroutine，避免协程泄漏
