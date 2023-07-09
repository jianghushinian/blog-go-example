package abs_test

import (
	"testing"

	"github.com/jianghushinian/blog-go-example/test/getting-started/abs"
)

// --------------黑盒测试---------------

func TestAbs(t *testing.T) {
	got := abs.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
