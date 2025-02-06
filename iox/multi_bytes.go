package iox

import (
	"io"
)

type MultiBytes struct {
	data  [][]byte // 存储数据的嵌套切片
	index int      // 当前读/写到的外层切片索引，data[index]
	pos   int      // 当前读/写到的切片所处理到的位置下标，data[index][pos]
}

var _ io.Reader = (*MultiBytes)(nil)
var _ io.Writer = (*MultiBytes)(nil)

// NewMultiBytes 构造一个 MultiBytes
func NewMultiBytes(data [][]byte) *MultiBytes {
	return &MultiBytes{
		data: data,
	}
}

// Read 实现 io.Reader 接口，从 data 中读取数据到 p
func (b *MultiBytes) Read(p []byte) (int, error) {
	// 如果 p 是空的，直接返回
	if len(p) == 0 {
		return 0, nil
	}

	// 所有数据都已读完
	if b.index >= len(b.data) {
		return 0, io.EOF
	}

	n := 0 // 记录已读取的字节数

	for n < len(p) {
		// 如果当前切片已经读完，则切换到下一个切片
		if b.pos >= len(b.data[b.index]) {
			b.index++
			b.pos = 0
			// 如果所有切片都已读完，退出循环
			if b.index >= len(b.data) {
				break
			}
		}

		// 从当前切片读取数据
		bytes := b.data[b.index]
		cnt := copy(p[n:], bytes[b.pos:])
		b.pos += cnt
		n += cnt
	}

	// 未读取到数据且已经读到结尾
	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

// Write 实现 io.Writer 接口，将数据追加到 data 中
func (b *MultiBytes) Write(p []byte) (int, error) {
	// 如果 p 是空的，直接返回
	if len(p) == 0 {
		return 0, nil
	}

	// 创建副本以避免外部修改影响数据
	clone := make([]byte, len(p))
	copy(clone, p)
	b.data = append(b.data, clone)
	return len(p), nil
}
