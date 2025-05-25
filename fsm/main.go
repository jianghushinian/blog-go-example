package main

import (
	"fmt"
)

type State string

const (
	ClosedState State = "closed"
	OpenState   State = "open"
)

type Event string

const (
	OpenEvent  Event = "open"
	CloseEvent Event = "close"
)

type Door struct {
	to    string
	state State
}

func NewDoor(to string) *Door {
	return &Door{
		to:    to,
		state: ClosedState,
	}
}

func (d *Door) CurrentState() State {
	return d.state
}

func (d *Door) HandleEvent(e Event) {
	switch e {
	case OpenEvent:
		d.state = OpenState
	case CloseEvent:
		d.state = ClosedState
	}
}

func main() {
	door := NewDoor("heaven")

	fmt.Println(door.CurrentState())

	door.HandleEvent(OpenEvent)
	fmt.Println(door.CurrentState())

	door.HandleEvent(CloseEvent)
	fmt.Println(door.CurrentState())
}
