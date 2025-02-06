package iox

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

// 测试创建一个 MultiBytes
func TestNewMultiBytes(t *testing.T) {
	// 创建参数
	type args struct {
		data [][]byte
	}

	// 测试用例
	tests := []struct {
		name string
		args args
		want *MultiBytes
	}{
		{
			name: "multi inner slice",
			args: args{data: [][]byte{
				[]byte("Hello, "),
				[]byte("世界！"),
				[]byte("This is an example of [][]byte to io.Reader and io.Writer."),
			}},
			want: &MultiBytes{
				data: [][]byte{
					[]byte("Hello, "),
					[]byte("世界！"),
					[]byte("This is an example of [][]byte to io.Reader and io.Writer."),
				},
			},
		},
		{
			name: "one inner slice",
			args: args{data: [][]byte{
				[]byte("Hello, World!"),
			}},
			want: &MultiBytes{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
		},
		{
			name: "outer slice is empty",
			args: args{data: [][]byte{}},
			want: &MultiBytes{
				data: [][]byte{},
			},
		},
		{
			name: "inner slice is empty",
			args: args{data: [][]byte{[]byte{}}},
			want: &MultiBytes{
				data: [][]byte{[]byte{}},
			},
		},
		{
			name: "outer slice is nil",
			args: args{data: nil},
			want: &MultiBytes{
				data: nil,
			},
		},
		{
			name: "inner slice is nil",
			args: args{data: [][]byte{nil}},
			want: &MultiBytes{
				data: [][]byte{nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMultiBytes(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMultiBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 MultiBytes 读功能
func TestMultiBytes_Read(t *testing.T) {
	// 创建 MultiBytes 参数
	type fields struct {
		data [][]byte
	}
	// 读取数据参数
	type args struct {
		p []byte
	}
	// 期待结果
	type want struct {
		n    int
		data []byte
		err  error
	}

	// 测试用例
	tests := []struct {
		name   string
		fields fields
		args   args
		before func(b *MultiBytes)
		after  func(b *MultiBytes)
		want   want
	}{
		{
			name: "multi inner slice read all",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, "),
					[]byte("世界！"),
					[]byte("This is an example of [][]byte to io.Reader and io.Writer."),
				},
			},
			args: args{p: make([]byte, 74)},
			want: want{
				n:    74,
				data: []byte("Hello, 世界！This is an example of [][]byte to io.Reader and io.Writer."),
			},
		},
		{
			name: "multi inner slice read half",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, "),
					[]byte("世界！"),
					[]byte("This is an example of [][]byte to io.Reader and io.Writer."),
				},
			},
			args: args{p: make([]byte, 34)},
			want: want{
				n:    34,
				data: []byte("Hello, 世界！This is an example"),
			},
		},
		{
			name: "multi inner slice read center",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, "),
					[]byte("世界！"),
					[]byte("This is an example of [][]byte to io.Reader and io.Writer."),
				},
			},
			args: args{p: make([]byte, 27)},
			before: func(b *MultiBytes) {
				// 先预读出前 7 个字节
				_, _ = b.Read(make([]byte, 7))
			},
			want: want{
				n:    27,
				data: []byte("世界！This is an example"),
			},
		},
		{
			name: "one inner slice read all",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
			args: args{p: make([]byte, 13, 100)},
			want: want{
				n:    13,
				data: []byte("Hello, World!"),
			},
		},
		{
			name: "one inner slice read half",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
			args: args{p: make([]byte, 5)},
			want: want{
				n:    5,
				data: []byte("Hello"),
			},
		},
		{
			name: "one inner slice read center",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
			args: args{p: make([]byte, 7)},
			before: func(b *MultiBytes) {
				// 先预读出前 5 个字节
				_, _ = b.Read(make([]byte, 5))
			},
			want: want{
				n:    7,
				data: []byte(", World"),
			},
		},
		{
			name: "outer slice is empty",
			fields: fields{
				data: [][]byte{},
			},
			args: args{p: make([]byte, 10)},
			want: want{
				n:    0,
				data: make([]byte, 10),
				err:  io.EOF,
			},
		},
		{
			name: "inner slice is empty",
			fields: fields{
				data: [][]byte{[]byte{}},
			},
			args: args{p: make([]byte, 10)},
			want: want{
				n:    0,
				data: make([]byte, 10),
				err:  io.EOF,
			},
		},
		{
			name: "outer slice is nil",
			fields: fields{
				data: nil,
			},
			args: args{p: make([]byte, 10)},
			want: want{
				n:    0,
				data: make([]byte, 10),
				err:  io.EOF,
			},
		},
		{
			name: "inner slice is nil",
			fields: fields{
				data: [][]byte{nil},
			},
			args: args{p: make([]byte, 10)},
			want: want{
				n:    0,
				data: make([]byte, 10),
				err:  io.EOF,
			},
		},
		{
			name: "p is nil",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
			args: args{p: nil},
			want: want{
				n:    0,
				data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewMultiBytes(tt.fields.data)

			if tt.before != nil {
				tt.before(b)
			}
			got, err := b.Read(tt.args.p)
			if tt.after != nil {
				tt.after(b)
			}

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.want.err)
			}
			if got != tt.want.n {
				t.Errorf("Read() got = %v, want %v", got, tt.want.n)
			}
			if !bytes.Equal(tt.args.p, tt.want.data) {
				t.Errorf("Read() got = %s, want %s", tt.args.p, tt.want.data)
			}
		})
	}
}

// 测试 MultiBytes 写功能
func TestMultiBytes_Write(t *testing.T) {
	// 创建 MultiBytes 参数
	type fields struct {
		data [][]byte
	}
	// 写入数据参数
	type args struct {
		p []byte
	}
	// 期待结果
	type want struct {
		n   int
		err error
	}

	p := []byte("Hello, World!")

	// 测试用例
	tests := []struct {
		name   string
		fields fields
		args   args
		before func(b *MultiBytes)
		after  func(b *MultiBytes)
		want   want
	}{
		{
			name: "normal",
			fields: fields{
				data: nil,
			},
			args: args{p: []byte("Hello, World!")},
			want: want{
				n: 13,
			},
		},
		{
			name: "data is not empty",
			fields: fields{
				data: [][]byte{
					[]byte("你好，世界！"),
				},
			},
			args: args{p: []byte("Hello, World!")},
			want: want{
				n: 13,
			},
		},
		{
			name: "p is nil",
			fields: fields{
				data: [][]byte{
					[]byte("Hello, World!"),
				},
			},
			args: args{p: nil},
			want: want{
				n: 0,
			},
		},
		{
			name: "modifying the given p does not affect the internal data",
			fields: fields{
				data: nil,
			},
			args: args{p: p},
			after: func(b *MultiBytes) { // FIXME: 这里也可以加入断言，验证是否正确
				// 写完成后修改原值 p
				p[5] = '.'
				t.Log(string(p)) // Hello. World!
				// 从 data 读取的内容不变
				newp := make([]byte, 13)
				_, _ = b.Read(newp)
				t.Log(string(newp)) // Hello, World!
			},
			want: want{
				n: 13,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewMultiBytes(tt.fields.data)

			if tt.before != nil {
				tt.before(b)
			}
			got, err := b.Write(tt.args.p)
			if tt.after != nil {
				tt.after(b)
			}

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.want.err)
			}
			if got != tt.want.n {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 MultiBytes 读写功能
func TestMultiBytes_ReadWrite(t *testing.T) {
	// TODO read-write-read
	// TODO write-read-write
}
