// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package bindingvalues

// Injectors from wire.go:

func InitializeMessage() Message {
	message := _wireMessageValue
	return message
}

var (
	_wireMessageValue = Message{
		Message: "Binding Values",
		Code:    1,
	}
)
