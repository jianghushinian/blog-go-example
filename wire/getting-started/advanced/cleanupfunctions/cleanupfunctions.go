package cleanupfunctions

import (
	"fmt"
	"os"
)

// NOTE: 清理函数

func OpenFile(path string) (*os.File, func(), error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		fmt.Println("cleanup...")
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}

	return f, cleanup, nil
}

func ReadFile(f *os.File) (string, error) {
	b := make([]byte, 1024)
	_, err := f.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
