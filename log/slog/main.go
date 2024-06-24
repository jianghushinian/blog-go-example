package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jianghushinian/blog-go-example/log/slog/customlog"
)

func main() {
	// 快速开始
	{
		slog.Debug("debug message")
		slog.Info("info message")
		slog.Warn("warn message")
		slog.Error("error message")
	}

	// 附加属性 key/value
	{
		slog.Debug("debug message", "hello", "world")
		slog.Info("info message", "hello", "world")
		slog.Warn("warn message", "hello", "world")
		slog.Error("error message", "hello", "world")
	}

	// context 版本
	{
		ctx := context.Background()
		slog.DebugContext(ctx, "debug message", "hello", "world")
		slog.InfoContext(ctx, "info message", "hello", "world")
		slog.WarnContext(ctx, "warn message", "hello", "world")
		slog.ErrorContext(ctx, "error message", "hello", "world")
	}

	// 修改日志级别
	{
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("debug message", "hello", "world")
		slog.Info("info message", "hello", "world")
		slog.Warn("warn message", "hello", "world")
		slog.Error("error message", "hello", "world")
	}

	// 获取当前日志级别
	// ref: https://stackoverflow.com/questions/77504588/how-to-retrieve-log-level-with-slog-package-in-go
	{
		var currentLevel slog.Level = -10
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			r := slog.Default().Enabled(context.Background(), level)
			if r {
				currentLevel = level
				break
			}
		}
		fmt.Printf("current log level: %v\n", currentLevel)
	}

	// 结构化日志
	// JSONHandler
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))
		l.Debug("debug message", "hello", "world")
	}
	// TextHandler
	{
		l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))
		l.Debug("debug message", "hello", "world")
	}

	// 使用自定义 logger 替换默认 logger
	{
		slog.Info("info message", "hello", "world")
		log.Println("normal log")

		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))
		slog.SetDefault(l)

		slog.Info("info message", "hello", "world")
		// log 也被修改了
		log.Println("normal log")
	}

	// 将 slog.Logger 转换为 log.Logger
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		logLogger := slog.NewLogLogger(l.Handler(), slog.LevelInfo)
		logLogger.Println("normal log") // 输出日志级别跟随 slog.LevelInfo 设置
	}

	// 使用宽松类型可能出现不匹配的 key/value
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		l.Info("info message", "hello") // "!BADKEY":"hello"

		// 使用 vet 命令检查此类问题
		// $ go vet .
		// ./main.go:120:3: call to slog.Logger.Info missing a final value
	}

	// 使用强类型
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		l.Info("info message", slog.String("hello", "world"), slog.Int("status", 200))
		l.Info("info message", slog.String("hello", "world"), slog.Int("status", 200), "extra") // "!BADKEY":"extra"
	}

	// 利用 LogAttrs 限制必须使用强类型，避免出现 `!BADKEY`
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		l.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"info message",
			slog.String("hello", "world"),
			slog.Int("status", 405),
			slog.Any("err", errors.New("http method not allowed")), // error 类型可以使用 slog.Any 输出
			// "extra","text", // 编译不通过，类型不匹配
		)
	}

	// 属性分组
	// JSONHandler
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		l.Info(
			"info message",
			slog.Group("user", "name", "root", slog.Int("age", 20)),
		)
	}
	// TextHandler
	{
		l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		l.Info(
			"info message",
			slog.Group("user", "name", "root", slog.Int("age", 20)),
		)
	}

	// 使用子 logger
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))
		// 附加自定义属性
		sl := l.With("requestId", "10191529-bc34-4efe-95e4-ecac7321773a")
		sl.Debug("debug message")
		sl.Info("info message")
	}

	// 为子 logger 属性分组
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		sl := l.WithGroup("user").With("requestId", "10191529-bc34-4efe-95e4-ecac7321773a")
		sl.Debug("debug message", "name", "admin")
		sl.Info("info message", "name", "admin")
	}

	// 实现 slog.LogValuer 接口，隐藏敏感信息
	{
		l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,            // 记录日志位置
			Level:       slog.LevelDebug, // 设置日志级别
			ReplaceAttr: nil,
		}))

		user := &User{
			ID:       123,
			Name:     "jianghushinian",
			Password: "pass",
		}
		l.Info("info message", "user1", user)  // *User 未实现 slog.LogValuer 接口
		l.Info("info message", "user2", *user) // User 未实现 slog.LogValuer 接口，所以无法隐藏敏感信息
		// l.Info("info message", "user3", user.SecureString())
	}

	// 使用自定义 logger
	{
		l := customlog.New(customlog.LevelDebug)
		l.Debug("custom debug message", "hello", "world")
		l.Trace("custom trace message", "hello", "world")
		l.Info("custom info message", "hello", "world")

		l.SetLevel(customlog.LevelInfo)
		l.Debug("custom debug message", "hello", "world")
		l.Trace("custom trace message", "hello", "world")
		l.Info("custom info message", "hello", "world")
	}

	// 使用自定义 handler
	{
		l := slog.New(customlog.NewHandler(os.Stdout, nil))
		l.Info("info message", "hello", "world")
	}
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LogValue implements slog.LogValuer interface
// slog.Value 不可比较: https://jianghushinian.cn/2024/06/15/how-to-make-structures-incomparable-in-go/
func (u *User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", u.ID),
		slog.String("name", u.Name),
	)
}

func (u *User) SecureString() string {
	u.Password = ""
	res, _ := json.Marshal(u)
	return string(res)
}
