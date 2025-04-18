package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func healthCheck() {
	// 发送健康检查请求
	resp, err := http.Get("http://127.0.0.1:5555/healthz")
	if err != nil {
		panic(fmt.Sprintf("请求失败: %v", err))
	}
	defer resp.Body.Close() // 确保连接回收

	// 丢弃响应体
	_, _ = io.Copy(io.Discard, resp.Body)

	// 状态码校验
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("非预期状态码: %d", resp.StatusCode))
	}

	fmt.Println("健康检查通过")
}

func filesize(name string) (int64, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, errors.New("文件不存在")
		}
		return 0, fmt.Errorf("读取文件失败: %w", err)
	}
	return fi.Size(), nil
}

func main() {
	healthCheck()
}
