package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jianghushinian/blog-go-example/design-patterns/builder/middleware/middleware/builder"
	"github.com/jianghushinian/blog-go-example/design-patterns/builder/middleware/middleware/options"
	"github.com/jianghushinian/blog-go-example/design-patterns/builder/middleware/middleware/simplebuilder"
)

func main() {
	r := gin.Default()

	// builder 模式中间件
	{
		// 为 /admin 路由设置只有管理员角色才能访问
		adminMiddleware := builder.NewRBACMiddlewareBuilder().
			AllowRole("admin").
			Build().
			Middleware()

		r.GET("/admin", adminMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome Admin!",
			})
		})

		// 为 /user 路由设置普通用户和管理员角色都能访问
		userMiddleware := builder.NewRBACMiddlewareBuilder().
			AllowRole("admin").
			AllowRole("user").
			Build().
			Middleware()

		r.GET("/user", userMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome User!",
			})
		})
	}

	// 极简版 builder 模式中间件
	{
		adminMiddleware := simplebuilder.NewRBACMiddlewareBuilder().
			AllowRole("admin").
			Build()

		r.GET("/simplebuilder/admin", adminMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome Admin!",
			})
		})

		userMiddleware := simplebuilder.NewRBACMiddlewareBuilder().
			AllowRole("admin").
			AllowRole("user").
			Build()

		r.GET("/simplebuilder/user", userMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome User!",
			})
		})
	}

	// options 模式中间件
	{
		adminMiddleware := options.NewRBACMiddleware(
			options.WithRole("admin"),
		).Middleware()

		r.GET("/options/admin", adminMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome Admin!",
			})
		})

		userMiddleware := options.NewRBACMiddleware(
			options.WithRole("user"),
			options.WithRole("admin"),
		).Middleware()

		r.GET("/options/user", userMiddleware, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome User!",
			})
		})
	}

	r.Run(":8000")
}
