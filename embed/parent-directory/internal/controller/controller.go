package controller

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"

	templatefs "github.com/jianghushinian/blog-go-example/embed/parent-directory/template"
)

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
