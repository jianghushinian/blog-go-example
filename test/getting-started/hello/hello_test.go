package hello_test

import (
	"testing"

	"github.com/jianghushinian/blog-go-example/test/getting-started/hello"
)

// --------------黑盒测试---------------

func TestHello_BlackBox(t *testing.T) {
	actual, err := hello.Hello("Tim")
	if err != nil {
		t.Error(err)
	}
	expected := "Hello Tim"
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}
}
