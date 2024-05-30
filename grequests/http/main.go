package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NOTE: 为了保持代码逻辑清晰，忽略了所有错误处理

func main() {
	// 使用 net/http 发送 GET 请求
	{
		// 发起 HTTP GET 请求
		resp, _ := http.Get("https://httpbin.org/get")
		defer resp.Body.Close() // 确保关闭响应体

		// 读取响应体内容
		body, _ := io.ReadAll(resp.Body)

		// 将响应体打印为字符串
		fmt.Println(string(body))

		// {
		//  "args": {},
		//  "headers": {
		//    "Accept-Encoding": "gzip",
		//    "Host": "httpbin.org",
		//    "User-Agent": "Go-http-client/2.0",
		//    "X-Amzn-Trace-Id": "Root=1-6657c816-1d63b01102c4d47e0dc10b2c"
		//  },
		//  "origin": "69.28.52.250",
		//  "url": "https://httpbin.org/get"
		// }
	}

	// 使用 net/http 发送 POST 请求
	{
		// 创建一个要发送的数据结构并编码为 JSON
		data := map[string]any{
			"username": "user",
			"password": "pass",
		}
		jsonData, _ := json.Marshal(data)

		// 创建 POST 请求
		url := "https://httpbin.org/post"
		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

		// 添加请求头，指明发送的内容类型为 JSON
		request.Header.Set("Content-Type", "application/json")

		// 发送请求并获取响应
		client := &http.Client{}
		resp, _ := client.Do(request)
		defer resp.Body.Close()

		// 读取响应体内容
		body, _ := io.ReadAll(resp.Body)

		// 将响应体打印为字符串
		fmt.Println(string(body))

		// {
		//  "args": {},
		//  "data": "{\"password\":\"pass\",\"username\":\"user\"}",
		//  "files": {},
		//  "form": {},
		//  "headers": {
		//    "Accept-Encoding": "gzip",
		//    "Content-Length": "37",
		//    "Content-Type": "application/json",
		//    "Host": "httpbin.org",
		//    "User-Agent": "Go-http-client/2.0",
		//    "X-Amzn-Trace-Id": "Root=1-6657c816-0743afe10cfc36a57ed145fc"
		//  },
		//  "json": {
		//    "password": "pass",
		//    "username": "user"
		//  },
		//  "origin": "69.28.52.250",
		//  "url": "https://httpbin.org/post"
		// }
	}
}
