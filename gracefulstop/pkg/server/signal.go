package server

import (
	"context"
	"os"
	"os/signal"
)

// https://github.com/kubernetes/apiserver/blob/release-1.31/pkg/server/signal.go

// 用来确保 SetupSignalContext 或 SetupSignalHandler 只被调用一次。如果尝试第二次调用，会触发 panic
var onlyOneSignalHandler = make(chan struct{})

// 用于接收操作系统信号的 channel，当接收到信号时，会通知相关的处理函数
var shutdownHandler chan os.Signal

// SetupSignalHandler 注册信号 SIGTERM 和 SIGINT，返回一个 channel
// 捕获到第一个信号进行优雅退出，捕获到第二个信号直接退出，退出码为 1
// SetupSignalContext 和 SetupSignalHandler 只能调用一个，并且只能调用一次
func SetupSignalHandler() <-chan struct{} {
	return SetupSignalContext().Done()
}

// SetupSignalContext 此函数与 SetupSignalHandler 函数区别是返回一个 context.Context
// SetupSignalContext 和 SetupSignalHandler 只能调用一个，并且只能调用一次
func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler) // 调用两次将 panic

	// 这里要注意：signal.Notify(c chan<- os.Signal, sig ...os.Signal) 函数不会为了向 c 发送信息而阻塞。
	// 也就是说，如果发送时 c 阻塞了，signal 包会直接丢弃信号。为了不丢失信号，我们创建了有缓冲的 channel  shutdownHandler。
	shutdownHandler = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler // 接收到第一个信号，走优雅退出流程
		cancel()          // 因为 ctx 被返回给调用方，取消 context 会通知到调用方
		<-shutdownHandler
		os.Exit(1) // 接收到第二个信号，直接退出
	}()

	return ctx
}

// RequestShutdown 主动触发优雅退出事件信号（SIGTERM/SIGINT），返回值表示是否触发成功
func RequestShutdown() bool {
	if shutdownHandler != nil {
		select {
		case shutdownHandler <- shutdownSignals[0]:
			return true
		default:
		}
	}

	return false
}
