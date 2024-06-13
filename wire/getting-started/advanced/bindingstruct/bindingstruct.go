package bindingstruct

import "fmt"

// NOTE: 绑定结构体到接口

type Message struct {
	Content string
	Code    int
}

type Store interface {
	Save(msg *Message) error
}

type store struct{}

// 确保 store 实现了 Store 接口
var _ Store = (*store)(nil)

func New() *store {
	return &store{}
}

func (s *store) Save(msg *Message) error {
	return nil
}

func SaveMessage(s Store, msg *Message) error {
	fmt.Printf("save message: %+v\n", msg)
	return s.Save(msg)
}

func RunStore(msg *Message) error {
	s := New()
	return SaveMessage(s, msg)
}
