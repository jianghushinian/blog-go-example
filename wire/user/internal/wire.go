//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/jianghushinian/blog-go-example/wire/user/internal/biz"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/config"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/controller"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/store"
	"github.com/jianghushinian/blog-go-example/wire/user/pkg/db"
)

// 依赖注入构造 Web 应用声明函数
func wireApp(engine *gin.Engine, cfg *config.Config, mysqlOptions *db.MySQLOptions) (*App, func(), error) {
	wire.Build(
		db.NewMySQL,
		store.ProviderSet,
		biz.ProviderSet,
		controller.New,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
