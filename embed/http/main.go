package main

import (
	"embed"
	"io/fs"
	"net/http"
	"sync"
	"text/template"
)

//go:embed static
var staticFS embed.FS

//go:embed template
var templateFS embed.FS

func main() {
	var wg sync.WaitGroup
	wg.Add(5)

	// NOTE: 在 go:embed 出现之前托管静态文件服务的写法
	go func() {
		defer wg.Done()
		http.Handle("/", http.FileServer(http.Dir("static")))
		_ = http.ListenAndServe(":8000", nil)
	}()

	// NOTE: 使用 go:embed 实现静态文件服务
	go func() {
		defer wg.Done()
		_ = http.ListenAndServe(":8001", http.FileServer(http.FS(staticFS)))
	}()

	// NOTE: 可以使用 http.StripPrefix 去除静态文件服务的 `/static/` 路由前缀
	go func() {
		defer wg.Done()
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
		_ = http.ListenAndServe(":8002", nil)
	}()

	// NOTE: 也可以使用 fs.Sub 去除静态文件服务的 `/static/` 路由前缀
	go func() {
		defer wg.Done()
		fsSub, _ := fs.Sub(staticFS, "static")
		_ = http.ListenAndServe(":8003", http.FileServer(http.FS(fsSub)))
	}()

	// NOTE: text/template 和 html/template 同样可以从嵌入的文件系统中解析模板，这里以 text/template 为例
	go func() {
		defer wg.Done()
		tmpl, _ := template.ParseFS(templateFS, "template/email.tmpl")
		http.HandleFunc("/email", func(writer http.ResponseWriter, request *http.Request) {
			// 设置 Content-Type 为 text/html
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")

			// 执行模板并发送响应
			_ = tmpl.ExecuteTemplate(writer, "email.tmpl", map[string]string{
				"Username":   "江湖十年",
				"ConfirmURL": "https://jianghushinian.cn",
			})
		})
		_ = http.ListenAndServe(":8004", nil)
	}()

	wg.Wait()
}
