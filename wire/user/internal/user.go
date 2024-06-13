package user

import (
	"github.com/gin-gonic/gin"

	"github.com/jianghushinian/blog-go-example/wire/user/internal/config"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/controller"
)

// App 代表一个 Web 应用
type App struct {
	*config.Config

	g  *gin.Engine
	uc *controller.UserController
}

// NewApp Web 应用构造函数
func NewApp(cfg *config.Config) (*App, func(), error) {
	// NOTE: 手动构造 app
	// gormDB, cleanup, err := db.NewMySQL(&cfg.MySQL)
	// if err != nil {
	// 	return nil, nil, err
	// }
	//
	// userStore := store.New(gormDB)
	// userBiz := biz.New(userStore)
	// userController := controller.New(userBiz)
	//
	// engine := gin.Default()
	// app := &App{
	// 	Config: cfg,
	// 	g:      engine,
	// 	uc:     userController,
	// }

	// NOTE: 依赖注入
	engine := gin.Default()
	app, cleanup, err := wireApp(engine, cfg, &cfg.MySQL)

	return app, cleanup, err
}

// Run 启动 Web 应用
func (a *App) Run() {
	// 注册路由
	InitRouter(a)

	if err := a.g.Run(":8000"); err != nil {
		panic(err)
	}
}
