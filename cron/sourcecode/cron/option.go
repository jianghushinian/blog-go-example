package cron

import (
	"time"
)

// Option 选项表示对 Cron 默认行为的修改。
type Option func(*Cron)

// WithLocation 覆盖 cron 实例的时区。
func WithLocation(loc *time.Location) Option {
	return func(c *Cron) {
		c.location = loc
	}
}

// WithSeconds 将覆盖用于解释作业执行计划的解析器，以将 seconds 字段作为第一个字段。
// 默认以 minutes 作为自一个字段
func WithSeconds() Option {
	return WithParser(NewParser(
		Second | Minute | Hour | Dom | Month | Dow | Descriptor,
	))
}

// WithParser 覆盖用于解释作业计划的解析器。
func WithParser(p ScheduleParser) Option {
	return func(c *Cron) {
		c.parser = p
	}
}

// WithChain 指定要应用于此 cron 的所有作业的 Job 装饰器。
func WithChain(wrappers ...JobWrapper) Option {
	return func(c *Cron) {
		c.chain = NewChain(wrappers...)
	}
}

// WithLogger 使用提供的日志记录器。
func WithLogger(logger Logger) Option {
	return func(c *Cron) {
		c.logger = logger
	}
}
