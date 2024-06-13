//go:build wireinject

package structfields

import "github.com/google/wire"

func InitializeMessage(phrase string, code int) Content {
	// 因为 NewMessage 返回 *Message 类型，而非 Message 类型，所以必须使用 new(*Message) 而不是 new(Message)
	// 此外，inject 函数参数和返回值类型都必须唯一，不然 wire 无法对应上哪个值是给谁的，所以定义了 Content 类型用来区分 message string
	wire.Build(NewMessage, wire.FieldsOf(new(*Message), "Content"))
	return Content("")
}
