//go:build wireinject

package providersets

import "github.com/google/wire"

func InitializeEvent(phrase string) (Event, error) {
	wire.Build(NewEvent, providerSet) // 顺序不重要
	return Event{}, nil
}
