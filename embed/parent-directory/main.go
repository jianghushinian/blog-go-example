package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jianghushinian/blog-go-example/embed/parent-directory/internal/controller"
)

func main() {
	r := gin.Default()
	r.POST("/users", (&controller.Controller{}).CreateUser)
	_ = r.Run(":8005")
}
