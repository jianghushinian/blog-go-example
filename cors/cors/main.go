package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NOTE: 在 Handler Func 中处理简单请求
// func main() {
// 	r := gin.Default()
//
// 	// 路由处理：返回简单的 JSON 数据
// 	r.GET("/data", func(c *gin.Context) {
// 		// 设置 CORS 响应头
// 		c.Header("Access-Control-Allow-Origin", "*")
//
// 		c.JSON(200, gin.H{
// 			"message": "这是配置了 CORS 的响应",
// 		})
// 	})
//
// 	// 监听在 8000 端口
// 	r.Run(":8000")
// }

// NOTE: 使用 CORS 中间件处理非简单请求
// func main() {
// 	r := gin.Default()
//
// 	// 使用 CORS 中间件
// 	r.Use(Cors)
//
// 	// 路由处理：返回简单的 JSON 数据
// 	r.POST("/data", func(c *gin.Context) {
// 		var requestData map[string]interface{}
// 		if err := c.BindJSON(&requestData); err != nil {
// 			c.JSON(400, gin.H{"error": "Invalid JSON"})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"message":     "非简单请求已成功",
// 			"requestData": requestData,
// 		})
// 	})
//
// 	// 监听在 8000 端口
// 	r.Run(":8000")
// }
//
// func Cors(c *gin.Context) {
// 	// 如果没有 Origin 请求头则说明不是 CORS 请求
// 	if c.Request.Header.Get("Origin") == "" {
// 		return
// 	}
//
// 	// 允许的 CORS 请求源
// 	c.Header("Access-Control-Allow-Origin", "*")
//
// 	// 处理预请求
// 	if c.Request.Method == "OPTIONS" {
// 		c.Header("Access-Control-Allow-Headers", "Content-Type,X-Custom-Header")
// 		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
// 		c.Header("Content-Type", "application/json")
// 		c.AbortWithStatus(204)
// 	}
//
// 	// 处理正式请求
// 	c.Next()
// }

// NOTE: 使用 Gin 的 CORS 中间件处理简单请求
// func main() {
// 	r := gin.Default()
//
// 	// 使用 Gin 的 CORS 中间件
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"*"}, // 允许的跨域源，"*" 表示任意源
// 		AllowMethods:     []string{"GET"},
// 		AllowHeaders:     []string{"Origin", "Content-Type"},
// 		AllowCredentials: false,
// 	}))
//
// 	// 路由处理：返回简单的 JSON 数据
// 	r.GET("/data", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "这是配置了 CORS 的响应",
// 		})
// 	})
//
// 	// 监听在 8000 端口
// 	r.Run(":8000")
// }

// NOTE: 使用 Gin 的 CORS 中间件处理非简单请求
func main() {
	r := gin.Default()

	// 使用 CORS 中间件配置非简单请求
	r.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:63342"},          // 允许的前端域名
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},          // 允许的 HTTP 方法
		AllowHeaders:     []string{"Content-Type", "X-Custom-Header"}, // 允许的请求头
		ExposeHeaders:    []string{"X-jwt-token"},                     // 允许暴露给 JavaScript 脚本的响应头
		AllowCredentials: true,                                        // 是否允许凭证（Cookies）
		AllowOriginFunc: func(origin string) bool { // 使用函数设置允许的前端域名
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.HasSuffix(origin, "jianghushinian.cn")
		},
		MaxAge: 1 * time.Hour,
	}))

	// 路由处理：返回简单的 JSON 数据
	r.POST("/data", func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}
		c.Header("X-jwt-token", "fake-token")
		c.JSON(200, gin.H{
			"message":     "非简单请求已成功",
			"requestData": requestData,
		})
	})

	// 监听在 8000 端口
	r.Run(":8000")
}
