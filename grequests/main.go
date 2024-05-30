package main

import (
	"fmt"
	"log"

	"github.com/levigross/grequests"
)

// NOTE: 为了保持代码逻辑清晰，忽略了所有错误处理

func main() {
	// GET
	{
		resp, _ := grequests.Get("https://httpbin.org/get", nil)
		defer resp.Close() // 确保关闭响应体

		fmt.Println(resp)
		fmt.Println(resp.String()) // 实现了 `fmt.Stringer` 接口
		fmt.Println(resp.RawResponse)
	}

	// GET query string
	{
		// 使用 RequestOptions 传递参数
		ro := &grequests.RequestOptions{
			Params: map[string]string{"Hello": "Goodbye"},
		}
		resp, _ := grequests.Get("https://httpbin.org/get?Hello=World", ro)
		defer resp.Close()

		fmt.Println(resp)
		// 现在 `url` 已经被替换成了 `https://httpbin.org/get?Hello=Goodbye`，可见 `ro` 会覆盖 `url` 中的同名参数
		fmt.Println(resp.RawResponse.Request.URL)
	}

	// GET set `User-Agent`
	{
		ro := &grequests.RequestOptions{UserAgent: "jianghushinian/bot 0.1"}
		resp, _ := grequests.Get("https://httpbin.org/get", ro)
		defer resp.Close()

		if resp.Ok != true {
			log.Println("Request did not return OK")
		}

		fmt.Println(resp.String())
	}

	// POST
	{
		// 创建一个要发送的数据结构
		postData := map[string]string{
			"username": "user",
			"password": "pass",
		}

		// 将数据结构编码为 JSON 并准备请求选项
		ro := &grequests.RequestOptions{
			JSON: postData, // grequests 自动处理 JSON 编码
			// Data: postData,
		}

		// 发起 POST 请求
		resp, _ := grequests.Post("https://httpbin.org/post", ro)
		defer resp.Close()

		// 输出响应体内容
		fmt.Println("Response:", resp.String())
	}

	// GET basic auth
	{
		ro := &grequests.RequestOptions{Auth: []string{"user", "pass"}}
		resp, _ := grequests.Get("https://httpbin.org/basic-auth/user/pass", ro)
		defer resp.Close()

		if resp.Ok != true {
			log.Println("Request did not return OK")
		}

		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header.Get("content-type"))
		fmt.Println(resp.String())

		m := make(map[string]any)
		_ = resp.JSON(&m)
		fmt.Println(m)
	}

	// GET download file
	{
		resp, _ := grequests.Get("https://httpbin.org/get", nil)
		defer resp.Close()

		if resp.Ok != true {
			log.Println("Request did not return OK")
		}

		// 下载响应体内容到 result.json
		_ = resp.DownloadToFile("result.json")
	}

	// POST upload file
	{
		// 从 result.json 读取文件内容并上传到服务端
		fd, _ := grequests.FileUploadFromDisk("result.json")

		// This will upload the file as a multipart mime request
		resp, _ := grequests.Post("https://httpbin.org/post",
			&grequests.RequestOptions{
				Files: fd,
				Data:  map[string]string{"One": "Two"},
			})
		defer resp.Close()

		if resp.Ok != true {
			log.Println("Request did not return OK")
		}

		fmt.Println(resp)
	}
}

// {
//  "args": {},
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Host": "httpbin.org",
//    "User-Agent": "GRequests/0.10",
//    "X-Amzn-Trace-Id": "Root=1-6657cca6-1957f8766defcf1502173532"
//  },
//  "origin": "69.28.52.250",
//  "url": "http://httpbin.org/get"
// }
//
// {
//  "args": {},
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Host": "httpbin.org",
//    "User-Agent": "GRequests/0.10",
//    "X-Amzn-Trace-Id": "Root=1-6657cca6-1957f8766defcf1502173532"
//  },
//  "origin": "69.28.52.250",
//  "url": "http://httpbin.org/get"
// }
//
// &{200 OK 200 HTTP/1.1 1 1 map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[*] Connection:[keep-alive] Content-Length:[266] Content-Type:[application/json] Date:[Thu, 30 May 2024 00:47:34 GMT] Server:[gunicorn/19.9.0]] 0x1400002c080 266 [] false false map[] 0x1400015e360 <nil>}

// {
//  "args": {
//    "Hello": "Goodbye"
//  },
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Host": "httpbin.org",
//    "User-Agent": "GRequests/0.10",
//    "X-Amzn-Trace-Id": "Root=1-6657cca7-715b811a08dff4420f487570"
//  },
//  "origin": "69.28.52.250",
//  "url": "https://httpbin.org/get?Hello=Goodbye"
// }
//
// https://httpbin.org/get?Hello=Goodbye

// {
//  "args": {},
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Host": "httpbin.org",
//    "User-Agent": "jianghushinian/bot 0.1",
//    "X-Amzn-Trace-Id": "Root=1-6657cca7-087fb9453524212d532c6bbe"
//  },
//  "origin": "69.28.52.250",
//  "url": "https://httpbin.org/get"
// }

// Response: {
//  "args": {},
//  "data": "{\"password\":\"pass\",\"username\":\"user\"}",
//  "files": {},
//  "form": {},
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Content-Length": "37",
//    "Content-Type": "application/json",
//    "Host": "httpbin.org",
//    "User-Agent": "GRequests/0.10",
//    "X-Amzn-Trace-Id": "Root=1-6657cca7-494c4a4078e419cd20b9c48c"
//  },
//  "json": {
//    "password": "pass",
//    "username": "user"
//  },
//  "origin": "69.28.52.250",
//  "url": "https://httpbin.org/post"
// }

// 200
// application/json
// {
//  "authenticated": true,
//  "user": "user"
// }
//
// map[authenticated:true user:user]

// {
//  "args": {},
//  "data": "",
//  "files": {
//    "file": "{\n  \"args\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"GRequests/0.10\", \n    \"X-Amzn-Trace-Id\": \"Root=1-6657cca8-0011a54031fd6b3a406079d3\"\n  }, \n  \"origin\": \"69.28.52.250\", \n  \"url\": \"https://httpbin.org/get\"\n}\n"
//  },
//  "form": {
//    "One": "Two"
//  },
//  "headers": {
//    "Accept-Encoding": "gzip",
//    "Content-Length": "625",
//    "Content-Type": "multipart/form-data; boundary=c23ff989aba9ecce6f979ab6a94719d486b27180ba05cfb83dd1f482ee5c",
//    "Host": "httpbin.org",
//    "User-Agent": "GRequests/0.10",
//    "X-Amzn-Trace-Id": "Root=1-6657cca8-22fdcd2f50ec993b7d6d32d8"
//  },
//  "json": null,
//  "origin": "69.28.52.250",
//  "url": "https://httpbin.org/post"
// }
