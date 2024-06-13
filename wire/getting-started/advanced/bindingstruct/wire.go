//go:build wireinject

package bindingstruct

import "github.com/google/wire"

func WireRunStore(msg *Message) error {
	// new(Store) 接口无需使用指针
	wire.Build(SaveMessage, New, wire.Bind(new(Store), new(*store)))
	return nil
}
