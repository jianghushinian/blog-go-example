//go:build wireinject

package multi

import "github.com/google/wire"

// InitializeApp 初始化 App
func InitializeApp() (*App, func(), error) {
	wire.Build(
		NewDatabaseConnection,
		NewLogFile,
		wire.Struct(new(App), "DB", "Log"),
	)
	return &App{}, nil, nil
}
