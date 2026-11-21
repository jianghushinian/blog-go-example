# test/race

本文示例对应博客文章：[Go 竞态检测器到底在检测什么？从 data race 到 happens-before](https://jianghushinian.cn/2026/11/21/go-race-detector/)

在 `test/race` 目录下，可以逐个运行示例：

```bash
go run -race ./counter
go run -race ./counter/atomic
go run -race ./counter/mutex
go run -race ./map
go run -race ./map/mutex
go run -race ./slice
go run -race ./slice/mutex
go run -race ./loopvar
go run -race ./loopvar/fixed
go run -race ./once
go run -race ./once/fixed
go run -race ./waitgroup
go run -race ./waitgroup/fixed
go run -race ./mutex
go run -race ./mutex/fixed
go run -race ./buffer
go run -race ./buffer/fixed
go run -race ./racecondition
go run -race ./happensbefore
```

| 目录 | 说明 |
| --- | --- |
| `counter` | 共享计数器，未同步的读改写 |
| `counter/atomic` | 用 `atomic.Int64` 修复计数器 |
| `counter/mutex` | 用 `sync.Mutex` 修复计数器 |
| `map` | 并发写 map |
| `map/mutex` | 用锁保护 map |
| `slice` | 并发 append 共享 slice |
| `slice/mutex` | 用锁保护 slice append |
| `loopvar` | 循环外声明的变量被闭包捕获 |
| `loopvar/fixed` | 每轮迭代创建局部副本 |
| `once` | 未同步的延迟初始化 |
| `once/fixed` | 用 `sync.Once` 修复 |
| `waitgroup` | `WaitGroup.Add` 放在 goroutine 内 |
| `waitgroup/fixed` | 用 `WaitGroup.Go` 修复 |
| `mutex` | 值 receiver 复制了 mutex，锁不住原对象 |
| `mutex/fixed` | 指针 receiver 修复 |
| `buffer` | 多个 goroutine 共享 `bytes.Buffer` |
| `buffer/fixed` | 每个 goroutine 使用局部 buffer，再合并 |
| `racecondition` | 没有 data race、但有逻辑竞态 |
| `happensbefore` | 用 channel 建立 happens-before |

`counter`、`map`、`slice`、`loopvar`、`once`、`waitgroup`、`buffer` 是故意保留数据竞争的示例，运行时会打印 `WARNING: DATA RACE`，部分示例可能直接 fatal 或 panic。

`mutex` 目录是故意写错的，`go vet ./mutex` 会报告 `copylocks`；运行结果也不会是 1000。`racecondition` 用来演示 `-race` 不会报告的逻辑竞态。
