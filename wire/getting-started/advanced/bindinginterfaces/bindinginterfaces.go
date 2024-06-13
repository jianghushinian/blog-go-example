package bindinginterfaces

import (
	"fmt"
	"io"
)

// NOTE: 绑定接口

func Write(w io.Writer, value any) {
	n, err := fmt.Fprintln(w, value)
	fmt.Printf("n: %d, err: %v\n", n, err)
}
