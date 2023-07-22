package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jianghushinian/blog-go-example/test/testable/reader/mocks"
)

func TestGetChangeLog(t *testing.T) {
	expected := ChangeLogSpec{
		Version: "v0.1.1",
		ChangeLog: `
# Changelog
All notable changes to this project will be documented in this file.
`,
	}

	f, err := os.CreateTemp("", "TEST_CHANGELOG")
	assert.NoError(t, err)
	defer func() {
		_ = f.Close()
		_ = os.RemoveAll(f.Name())
	}()

	data := `
# Changelog
All notable changes to this project will be documented in this file.
`
	_, err = f.WriteString(data)
	assert.NoError(t, err)
	_, _ = f.Seek(0, 0)

	actual, err := GetChangeLog(f)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

type fakeReader struct {
	data   string
	offset int
}

func NewFakeReader(input string) io.Reader {
	return &fakeReader{
		data:   input,
		offset: 0,
	}
}

func (r *fakeReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF // 表示数据已读取完毕
	}

	n := copy(p, r.data[r.offset:]) // 将数据从字符串复制到 p 中
	r.offset += n

	return n, nil
}

func TestGetChangeLogByIOReader(t *testing.T) {
	expected := ChangeLogSpec{
		Version: "v0.1.1",
		ChangeLog: `
# Changelog
All notable changes to this project will be documented in this file.
`,
	}

	data := `
# Changelog
All notable changes to this project will be documented in this file.
`
	reader := NewFakeReader(data)
	actual, err := GetChangeLogByIOReader(reader)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetChangeLogByIOReader_mock(t *testing.T) {
	expected := ChangeLogSpec{
		Version: "v0.1.1",
		ChangeLog: `
# Changelog
All notable changes to this project will be documented in this file.
`,
	}

	data := `
# Changelog
All notable changes to this project will be documented in this file.
`
	ctrl := gomock.NewController(t)
	// reader := mocks.NewMockReaderWrapper(ctrl)
	reader := mocks.NewMockIReader(ctrl)
	reader.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		copy(p, data)
		return len(data), io.EOF
	})

	actual, err := GetChangeLogByIOReader(reader)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func init() {
	version = "v0.1.1"
}
