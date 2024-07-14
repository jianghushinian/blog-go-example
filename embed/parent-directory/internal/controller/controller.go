package controller

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"

	templatefs "github.com/jianghushinian/blog-go-example/embed/parent-directory/template"
)

// 直接嵌入父目录会报错
// //go:embed ../../template
// internal/controller/controller.go:12:12: pattern ../../template: invalid pattern syntax

// 不好的解决方案: ref https://github.com/golang/go/issues/46056
// 1. 相当于存储了两份依赖文件
// 2. 容易忘记执行 go generate ./...
// https://github.com/golang/go/issues/46056#issuecomment-938251415
//
// //go:generate cp -r ../../template/ template
// //go:embed template
// var templateFS embed.FS

type Controller struct{}

func (ctrl *Controller) CreateUser(c *gin.Context) {
	// pretend to create user ...

	// send email
	// tmpl, _ := template.ParseFS(templateFS, "template/*.tmpl")
	tmpl, _ := template.ParseFS(templatefs.TemplateFS, "*.tmpl") // 记得去掉 template 前缀
	var buf bytes.Buffer
	_ = tmpl.ExecuteTemplate(&buf, "email.tmpl", map[string]string{
		"Username":   "江湖十年",
		"ConfirmURL": "https://jianghushinian.cn",
	})
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
