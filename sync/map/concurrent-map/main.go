package main

import (
	"fmt"
)

// NOTE: Go 内置 map 不支持并发读写

func main() {
	m := make(map[string]int)

	// 并发读写 map
	go func() {
		for {
			m["k"] = 1
			fmt.Println("set k:", 1)
		}
	}()

	go func() {
		for {
			v, _ := m["k"]
			fmt.Println("read k:", v)
		}
	}()

	select {}
}

/*
$ go run main.go
fatal error: concurrent map read and map write
*/

/*
# 测试并发读写场景性能
$ go test -bench="Map$" -benchmem -run="^$"
goos: darwin
goarch: amd64
pkg: github.com/jianghushinian/blog-go-example/sync/concurrent-map
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkChannelMap-4             551652              2670 ns/op             478 B/op          6 allocs/op
BenchmarkRWMutexMap-4            1000000              1474 ns/op             196 B/op          4 allocs/op
BenchmarkMutexMap-4              2156851               571.5 ns/op            68 B/op          4 allocs/op
BenchmarkShardingMap-4           2752004               444.5 ns/op            66 B/op          3 allocs/op
BenchmarkSyncMap-4               1765488               618.8 ns/op           115 B/op          7 allocs/op
PASS
ok      github.com/jianghushinian/blog-go-example/sync/concurrent-map   10.271s
*/

/*
# 测试读多写少场景性能
$ go test -bench="MapReadHeavy$" -benchmem -run="^$"
goos: darwin
goarch: amd64
pkg: github.com/jianghushinian/blog-go-example/sync/concurrent-map
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkChannelMapReadHeavy-4              8490            212698 ns/op           33944 B/op        489 allocs/op
BenchmarkRWMutexMapReadHeavy-4             30123             37994 ns/op            5342 B/op        230 allocs/op
BenchmarkMutexMapReadHeavy-4               33726             42250 ns/op            5352 B/op        230 allocs/op
BenchmarkShardingMapReadHeavy-4            31478             36449 ns/op            5320 B/op        230 allocs/op
BenchmarkSyncMapReadHeavy-4                35464             36445 ns/op            5640 B/op        250 allocs/op
PASS
ok      github.com/jianghushinian/blog-go-example/sync/concurrent-map   10.816s
*/
