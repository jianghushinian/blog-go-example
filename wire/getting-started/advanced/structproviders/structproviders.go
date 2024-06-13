package structproviders

import (
	"errors"
	"fmt"
	"time"
)

// NOTE: 结构体 Provider

type Message struct {
	Content string
	Code    int
}

// 假设不提供 Message 构造函数
// func NewMessage(content string, code int) Message {
// 	return Message{
// 		Content: content,
// 		Code:    code,
// 	}
// }

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) (Event, error) {
	// 模拟创建 Event 报错
	if time.Now().Unix()%2 == 0 {
		return Event{}, errors.New("new event error")
	}
	return Event{Greeter: g}, nil
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Printf("%+v\n", msg)
}
