package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jianghushinian/blog-go-example/wire/user/internal/biz"
	"github.com/jianghushinian/blog-go-example/wire/user/pkg/api"
)

// UserController 用来处理用户请求
type UserController struct {
	b biz.UserBiz
}

// New controller 构造函数
func New(b biz.UserBiz) *UserController {
	return &UserController{b: b}
}

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	var r api.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	if err := ctrl.b.Create(c, &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
