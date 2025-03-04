package nightwatch

import (
	"log/slog"
)

type cronLogger struct{}

func newCronLogger() *cronLogger {
	return &cronLogger{}
}

func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	slog.Error(msg, append(keysAndValues, "err", err)...)
}
