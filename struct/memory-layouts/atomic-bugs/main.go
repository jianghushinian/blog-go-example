package main

import (
	"sync/atomic"
)

type S1 struct {
	a int32
	b int64
}

type S2 struct {
	a   int32
	pad uint32 // ensure 8-byte alignment of val on 386, ref: https://github.com/golang/go/issues/36606
	b   int64
}

func main() {
	{
		s1 := S1{}
		atomic.AddInt64(&s1.b, 1)
	}

	{
		s2 := S2{}
		atomic.AddInt64(&s2.b, 1) // 手动对齐以后不会报错
	}
}
