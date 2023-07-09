package hello

import (
	"errors"
	"testing"
)

// --------------模糊测试---------------

func FuzzHello(f *testing.F) {
	f.Add("Foo")
	f.Fuzz(func(t *testing.T, name string) {
		_, err := Hello(name)
		if err != nil {
			if errors.Is(err, ErrEmptyName) || errors.Is(err, ErrTooLongName) {
				return
			}
			t.Errorf("unexpected error: %s, name: %s", err, name)
		}
	})
}
