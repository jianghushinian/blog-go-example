//go:build wireinject

package bindingvalues

import (
	"github.com/google/wire"
)

func InitializeMessage() Message {
	// 假设没有提供 NewMessage，可以直接绑定值并返回
	wire.Build(wire.Value(Message{
		Message: "Binding Values",
		Code:    1,
	}))
	return Message{}
}
