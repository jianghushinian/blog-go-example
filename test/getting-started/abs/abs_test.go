package abs

import (
	"os"
	"testing"
	"time"
)

// --------------单元测试---------------

func TestAbs(t *testing.T) {
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}

func TestAbs_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{
			name: "positive",
			x:    2,
			want: 2,
		},
		{
			name: "negative",
			x:    -3,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%f) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestAbs_Skip(t *testing.T) {
	// CI 环境跳过当前测试
	if os.Getenv("CI") != "" {
		t.Skip("it's too slow, skip when running in CI")
	}

	t.Log(t.Skipped())

	got := Abs(-2)
	if got != 2 {
		t.Errorf("Abs(-2) = %f; want 2", got)
	}
}

func TestAbs_Parallel(t *testing.T) {
	t.Log("Parallel before")
	// 标记当前测试支持并行
	t.Parallel()
	t.Log("Parallel after")

	got := Abs(2)
	if got != 2 {
		t.Errorf("Abs(2) = %f; want 2", got)
	}
}

// 错误写法，不会被执行，被测试函数首字母应该大写
func Testabs(t *testing.T) {
	got := Abs(-2)
	if got != 2 {
		t.Errorf("Abs(-2) = %f; want 2", got)
	}
}

// --------------基准测试（性能测试）---------------

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

// 重置计时器
func BenchmarkAbsResetTimer(b *testing.B) {
	time.Sleep(100 * time.Millisecond) // 模拟耗时的准备工作
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

// 未重置计时器 测试结果不准确
func BenchmarkAbsExpensive(b *testing.B) {
	time.Sleep(100 * time.Millisecond) // 模拟耗时的准备工作
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

// 重置计时器另一种方式，先停止再开始
func BenchmarkAbsStopTimerStartTimer(b *testing.B) {
	b.StopTimer()
	time.Sleep(100 * time.Millisecond) // 模拟耗时的准备工作
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

// 并行测试
func BenchmarkAbsParallel(b *testing.B) {
	b.SetParallelism(2) // 设置并发 Goroutines 数量为 p * GOMAXPROCS
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Abs(-1)
		}
	})
}
