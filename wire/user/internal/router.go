package user

// InitRouter 初始化路由
func InitRouter(a *App) {
	// 创建 users 路由分组
	u := a.g.Group("/users")
	{
		u.POST("", a.uc.Create)
	}
}
