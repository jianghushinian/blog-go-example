//go:build wireinject

package argsanderror

import "github.com/google/wire"

func InitializeEvent(phrase string) (Event, error) {
	wire.Build(NewEvent, NewMessage, NewGreeter) // 顺序不重要
	return Event{}, nil
}
