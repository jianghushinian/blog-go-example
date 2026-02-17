package main

import (
	"embed"
	"io/fs"
	"net/http"
	"sync"
)

//go:embed static
var staticFS embed.FS

func main() {
	var wg sync.WaitGroup
	wg.Add(5)

	// 使用 go:embed 实现静态文件服务
	go func() {
		defer wg.Done()
		_ = http.ListenAndServe(":8001", http.FileServer(http.FS(staticFS)))
	}()

	// 使用 fs.Sub 去除静态文件服务的 `/static/` 路由前缀
	go func() {
		defer wg.Done()
		fsSub, _ := fs.Sub(staticFS, "static")
		_ = http.ListenAndServe(":8002", http.FileServer(http.FS(fsSub)))
	}()

	wg.Wait()
}
