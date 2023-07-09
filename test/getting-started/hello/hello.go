package hello

import "errors"

var (
	ErrEmptyName   = errors.New("empty name")
	ErrTooLongName = errors.New("too long name")
)

func Hello(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyName
	}
	if len(name) > 10 {
		return "", ErrTooLongName
	}
	// if name == "Bob" {
	// 	return "", errors.New("not allowed")
	// }
	return "Hello " + name, nil
}
