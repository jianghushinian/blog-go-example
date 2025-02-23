package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// 安装：go get github.com/robfig/cron/v3@master
// 一定要安装 master 分支版本，go get -u github.com/robfig/cron/v3 得到的 @v3.0.1 分支 cron.SkipIfStillRunning 方法存在 bug
// 当作业触发 panic 时，可能会永久阻塞后续作业的执行

// NOTE: 基础用法

// func main() {
// 	// 创建一个新的 Cron 实例
// 	c := cron.New(cron.WithSeconds())
//
// 	// 添加一个每秒执行的任务
// 	_, err := c.AddFunc("* * * * * *", func() {
// 		fmt.Println("每秒执行的任务:", time.Now().Format("2006-01-02 15:04:05"))
// 	})
// 	if err != nil {
// 		log.Fatalf("添加任务失败: %v", err)
// 	}
//
// 	// 添加一个每 5 秒执行的任务
// 	_, err = c.AddFunc("*/5 * * * * *", func() {
// 		fmt.Println("每 5 秒执行的任务:", time.Now().Format("2006-01-02 15:04:05"))
// 	})
// 	if err != nil {
// 		log.Fatalf("添加任务失败: %v", err)
// 	}
//
// 	// 启动 Cron
// 	c.Start()
// 	defer c.Stop() // 确保程序退出时停止 Cron
//
// 	// 主程序等待 10 秒，以便观察任务执行
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("主程序结束")
// }

// NOTE: 进阶用法

// 自定义 logger
type cronLogger struct{}

func newCronLogger() *cronLogger {
	return &cronLogger{}
}

// Info implements the cron.Logger interface.
func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

// Error implements the cron.Logger interface.
func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	slog.Error(msg, append(keysAndValues, "err", err)...)
}

// Job 作业对象
type Job struct {
	name  string
	count int
}

func (j *Job) Run() {
	j.count++
	if j.count == 2 {
		panic("第 2 次执行触发 panic")
	}
	if j.count == 4 {
		time.Sleep(6 * time.Second)
		log.Println("第 4 次执行耗时 6s")
	}
	fmt.Printf("%s: 每 5 秒执行的任务, count: %d\n", j.name, j.count)
}

func main() {
	log.Println("cron start")
	// 创建自定义日志对象
	logger := &cronLogger{}

	// 创建一个新的 Cron 实例
	c := cron.New(
		cron.WithSeconds(),      // 增加秒解析
		cron.WithLogger(logger), // 自定义日志
		cron.WithChain( // chain 是顺序敏感的
			cron.SkipIfStillRunning(logger), // 如果作业仍在运行，则跳过此次运行
			cron.Recover(logger),            // 恢复 panic
		),
	)

	var (
		spec = "@every 5s"
		job  = &Job{name: "江湖十年"}
	)

	// 添加一个每 5 秒执行的任务
	id, err := c.AddJob(spec, job)
	if err != nil {
		log.Fatalf("添加任务失败: %v", err)
	}
	log.Println("任务 ID:", id)

	// 启动 Cron
	c.Start()
	defer c.Stop() // 确保程序退出时停止 Cron

	time.Sleep(34 * time.Second) // 确保 job 能执行 6 次
	c.Remove(id)                 // 从调度器中移除 job
	time.Sleep(10 * time.Second) // job 不会再次执行
	log.Println("cron done")
}

func SkipIfStillRunning(logger cron.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		var ch = make(chan struct{}, 1)
		ch <- struct{}{}
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				defer func() { // 确保即使出现 panic，也不要阻塞后续作业的执行
					ch <- v
				}()
				j.Run()
			default:
				logger.Info("skip")
			}
		})
	}
}
