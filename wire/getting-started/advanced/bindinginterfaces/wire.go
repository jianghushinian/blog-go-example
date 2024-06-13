//go:build wireinject

package bindinginterfaces

import (
	"io"
	"os"

	"github.com/google/wire"
)

func InitializeWriter() io.Writer {
	wire.Build(wire.InterfaceValue(new(io.Writer), os.Stdout))
	return nil
}
