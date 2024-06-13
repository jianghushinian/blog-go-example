//go:build wireinject

package cleanupfunctions

import (
	"os"

	"github.com/google/wire"
)

func InitializeFile(path string) (*os.File, func(), error) {
	wire.Build(OpenFile)
	return nil, nil, nil
}
