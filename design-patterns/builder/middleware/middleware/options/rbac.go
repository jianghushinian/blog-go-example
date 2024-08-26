package options

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware RBAC 中间件结构体
type RBACMiddleware struct {
	allowedRoles []string
}

// Middleware 返回一个 Gin 中间件函数
func (r *RBACMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("Role") // 从请求头中获取用户角色
		for _, role := range r.allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
