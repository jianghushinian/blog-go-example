package builder

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

// RBACMiddlewareBuilder 用于构建 RBACMiddleware 的构建器
type RBACMiddlewareBuilder struct {
	rbacMiddleware *RBACMiddleware
}

// NewRBACMiddlewareBuilder 创建一个新的 RBACMiddlewareBuilder 实例
func NewRBACMiddlewareBuilder() *RBACMiddlewareBuilder {
	return &RBACMiddlewareBuilder{
		rbacMiddleware: &RBACMiddleware{},
	}
}

// AllowRole 添加允许访问的角色
func (b *RBACMiddlewareBuilder) AllowRole(role string) *RBACMiddlewareBuilder {
	b.rbacMiddleware.allowedRoles = append(b.rbacMiddleware.allowedRoles, role)
	return b
}

// Build 返回构建完成的 RBACMiddleware
func (b *RBACMiddlewareBuilder) Build() *RBACMiddleware {
	return b.rbacMiddleware
}
