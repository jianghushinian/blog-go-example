//go:build wireinject

package structproviders

import "github.com/google/wire"

func InitializeEvent(phrase string, code int) (Event, error) {
	// `*` 表示通配符，不显式指定的字段不会被赋值
	// wire.Build(NewEvent, NewGreeter, wire.Struct(new(Message), "*"))
	wire.Build(NewEvent, NewGreeter, wire.Struct(new(Message), "Content"))
	return Event{}, nil
}
