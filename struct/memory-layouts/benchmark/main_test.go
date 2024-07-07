package main

import "testing"

func Benchmark_T1_Align(b *testing.B) {
	b.ReportAllocs() // 开启内存统计
	b.ResetTimer()   // 重置计时器
	for i := 0; i < b.N; i++ {
		_ = make([]T1, b.N)
	}
}

func Benchmark_T2_Align(b *testing.B) {
	b.ReportAllocs() // 开启内存统计
	b.ResetTimer()   // 重置计时器
	for i := 0; i < b.N; i++ {
		_ = make([]T2, b.N)
	}
}
