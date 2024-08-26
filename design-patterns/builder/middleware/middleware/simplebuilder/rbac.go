package simplebuilder

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACMiddlewareBuilder RBAC 中间件结构体
type RBACMiddlewareBuilder struct {
	allowedRoles []string
}

// NewRBACMiddlewareBuilder 创建一个新的 RBACMiddlewareBuilder 实例
func NewRBACMiddlewareBuilder() *RBACMiddlewareBuilder {
	return &RBACMiddlewareBuilder{}
}

// Build 返回一个 Gin 中间件函数
func (b *RBACMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("Role") // 从请求头中获取用户角色
		for _, role := range b.allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

// AllowRole 添加允许访问的角色
func (b *RBACMiddlewareBuilder) AllowRole(role string) *RBACMiddlewareBuilder {
	b.allowedRoles = append(b.allowedRoles, role)
	return b
}
