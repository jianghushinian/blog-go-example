package main

import (
	"io"
	"os"
)

var version = "dev"

type ChangeLogSpec struct {
	Version   string
	ChangeLog string
}

// Bad: 依赖文件对象

func GetChangeLog(f *os.File) (ChangeLogSpec, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return ChangeLogSpec{}, err
	}

	return ChangeLogSpec{
		Version:   version,
		ChangeLog: string(data),
	}, nil
}

// Good: 使用接口解耦

func GetChangeLogByIOReader(reader io.Reader) (ChangeLogSpec, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return ChangeLogSpec{}, err
	}

	return ChangeLogSpec{
		Version:   version,
		ChangeLog: string(data),
	}, nil
}
