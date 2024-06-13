//go:build wireinject

package alternateinjectorsyntax

import "github.com/google/wire"

func InitializeMessage(phrase string) Message {
	panic(wire.Build(NewMessage))
}
