package customlog

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelTrace = slog.Level(-2) // 自定义日志级别
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Logger struct {
	l   *slog.Logger
	lvl *slog.LevelVar // 用来动态调整日志级别
}

func New(level slog.Level) *Logger {
	var lvl slog.LevelVar
	lvl.Set(level)

	h := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,

		// Level:     level, // 静态设置日志级别
		Level: &lvl, // 支持动态设置日志级别

		// 修改日志中的 Attr 键值对（即日志记录中附加的 key/value）
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel := level.String()

				switch level {
				case LevelTrace:
					// NOTE: 如果不设置，默认日志级别打印为 "level":"DEBUG+2"
					levelLabel = "TRACE"
				}

				a.Value = slog.StringValue(levelLabel)
			}

			// NOTE: 可以在这里修改时间输出格式
			// if a.Key == slog.TimeKey {
			// 	if t, ok := a.Value.Any().(time.Time); ok {
			// 		a.Value = slog.StringValue(t.Format(time.DateTime))
			// 	}
			// }

			return a
		},
	}))

	return &Logger{l: h, lvl: &lvl}
}

// SetLevel 动态调整日志级别
func (l *Logger) SetLevel(level Level) {
	l.lvl.Set(level)
}

func (l *Logger) Debug(msg string, args ...any) {
	// 不会走 *customlog.Logger.log() 调用，会走 *slog.Logger.log() 调用
	l.l.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	// l.l.Info(msg, args...)
	l.Log(context.Background(), LevelInfo, msg, args...)
}

// Trace 自定义的日志级别
func (l *Logger) Trace(msg string, args ...any) {
	l.Log(context.Background(), LevelTrace, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	// l.l.Warn(msg, args...)
	l.Log(context.Background(), LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	// l.l.Error(msg, args...)
	l.Log(context.Background(), LevelError, msg, args...)
}

func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.log(ctx, level, msg, args...)
}

// NOTE: 这个时候，可以自己设置 skip 大小，也可以自由决定是否启用 runtime.Callers
// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	// NOTE: 这里修改 skip 为 4，*slog.Logger.log 源码中 skip 为 3
	runtime.Callers(4, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.l.Handler().Handle(ctx, r)
}
