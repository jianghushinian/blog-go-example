package reader

import "io"

//go:generate mockgen -package mocks -destination mocks/ireader.go github.com/jianghushinian/blog-go-example/test/testable/reader IReader
//go:generate mockgen -package mocks -destination mocks/readerwrapper.go github.com/jianghushinian/blog-go-example/test/testable/reader ReaderWrapper

type ReaderWrapper interface {
	io.Reader
}

type IReader io.Reader
