package alternateinjectorsyntax

// NOTE: 使用 panic 简化 inject 代码

type Message string

func NewMessage(phrase string) Message {
	return Message(phrase)
}
