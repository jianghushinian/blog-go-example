// 模拟 Go module proxy,用于验证 GOPROXY 中逗号(,)与管道(|)分隔符的回退区别。
//
// PROXY_STATUS 决定该实例的行为:
//
//	200  -> 作为"健康代理",按 GOPROXY 协议动态提供任意模块的合法内容
//	        (list/info/mod/zip),让 go get 能真实下载成功;
//	404  -> 对一切请求返回 404,模拟"该代理上找不到模块";
//	500/502/503 -> 对一切请求返回对应状态码,模拟"代理服务端出错"。
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const modTime = "2024-01-01T00:00:00Z"

func main() {
	port := flag.String("port", "8080", "监听端口")
	status := flag.Int("status", 200, "模拟的 HTTP 状态码 (200/404/500/502/503)")
	name := flag.String("name", "proxy", "代理名称（用于日志）")
	flag.Parse()

	// 未传命令行参数时,从环境变量读取(方便 test.sh 与 README 使用)
	if len(os.Args) == 1 {
		if v := os.Getenv("PROXY_PORT"); v != "" {
			*port = v
		}
		if v := os.Getenv("PROXY_STATUS"); v != "" {
			fmt.Sscanf(v, "%d", status)
		}
		if v := os.Getenv("PROXY_NAME"); v != "" {
			*name = v
		}
	}

	h := &handler{name: *name, status: *status}
	addr := ":" + *port
	log.Printf("代理服务器启动: %s (名称: %s, 状态: %d)", addr, *name, *status)
	log.Fatal(http.ListenAndServe(addr, h))
}

type handler struct {
	name   string
	status int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s", h.name, r.Method, r.URL.Path)

	if h.status != 200 {
		// 统一模拟一个"坏代理":无论请求什么模块都返回同一状态码。
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Proxy-Name", h.name)
		w.WriteHeader(h.status)
		fmt.Fprintf(w, `{"proxy":%q,"simulated_status":%d}`, h.name, h.status)
		return
	}

	// 状态 200:按 GOPROXY 协议提供模块内容,模块路径从请求 URL 推导,
	// 因此无需硬编码,任意模块（含父路径 example.com 等）都能被合法提供。
	mod, suffix, ok := splitRequest(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respondModule(w, mod, suffix)
}

// splitRequest 把 /<module>/@v/<suffix> 拆成 module 与 suffix。
// module 中的大写字母在协议里以 "!小写" 形式传输,这里还原。
func splitRequest(path string) (mod, suffix string, ok bool) {
	p := strings.TrimPrefix(path, "/")
	if i := strings.LastIndex(p, "/@v/"); i >= 0 {
		return unescapeModule(p[:i]), p[i+len("/@v/"):], true
	}
	if i := strings.LastIndex(p, "/@latest"); i >= 0 {
		return unescapeModule(p[:i]), "@latest", true
	}
	return "", "", false
}

func unescapeModule(esc string) string {
	if !strings.Contains(esc, "!") {
		return esc
	}
	var b strings.Builder
	for i := 0; i < len(esc); i++ {
		if esc[i] == '!' && i+1 < len(esc) {
			c := esc[i+1]
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
			b.WriteByte(c)
			i++
		} else {
			b.WriteByte(esc[i])
		}
	}
	return b.String()
}

func respondModule(w http.ResponseWriter, mod, suffix string) {
	version := ""
	switch {
	case suffix == "@latest":
		version = "v1.2.0"
		writeInfo(w, mod, version)
	case suffix == "list":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "v1.0.0")
		fmt.Fprintln(w, "v1.1.0")
		fmt.Fprintln(w, "v1.2.0")
	case strings.HasSuffix(suffix, ".info"):
		version = strings.TrimSuffix(suffix, ".info")
		writeInfo(w, mod, version)
	case strings.HasSuffix(suffix, ".mod"):
		version = strings.TrimSuffix(suffix, ".mod")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, modFile(mod))
	case strings.HasSuffix(suffix, ".zip"):
		version = strings.TrimSuffix(suffix, ".zip")
		data, err := moduleZip(mod, version)
		if err != nil {
			log.Printf("[moduleZip] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Write(data)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func writeInfo(w http.ResponseWriter, mod, version string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Version": version,
		"Time":    modTime,
	})
}

func modFile(mod string) string {
	return fmt.Sprintf("module %s\n\ngo 1.21\n", mod)
}

// moduleZip 生成符合 GOPROXY 协议的合法 zip:
// 目录为 <module>@<version>/,内含与 .mod 端点一致的 go.mod 和一个包文件。
func moduleZip(mod, version string) ([]byte, error) {
	prefix := mod + "@" + version + "/"
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	add := func(name, content string) error {
		f, err := zw.Create(prefix + name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(content))
		return err
	}
	if err := add("go.mod", modFile(mod)); err != nil {
		return nil, err
	}
	if err := add(pkgName(mod)+".go", fmt.Sprintf("package %s\n", pkgName(mod))); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// pkgName 取模块路径最后一段作为包名,并去掉非法字符。
func pkgName(mod string) string {
	name := mod
	if i := strings.LastIndex(mod, "/"); i >= 0 {
		name = mod[i+1:]
	}
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, name)
	if name == "" {
		return "module"
	}
	return name
}
