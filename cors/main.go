package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路由处理：返回简单的 JSON 数据
	r.GET("/data", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "这是未配置 CORS 的响应",
		})
	})

	// 监听在 8000 端口
	r.Run(":8000")
}
