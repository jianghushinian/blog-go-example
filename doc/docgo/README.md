## 使用

```shell
# 启动 govanityurls 服务
# https://github.com/GoogleCloudPlatform/govanityurls
$ go install github.com/GoogleCloudPlatform/govanityurls
$ cd govanityurls/
$ govanityurls

# 启动 nginx 反向代理
$ nginx-https/
$ docker run -d --name nginx-https \
  -p 80:80 -p 443:443 \
  -v $(pwd)/conf.d:/etc/nginx/conf.d \
  -v $(pwd)/ssl:/etc/nginx/ssl \
  nginx:1.23.3

$ vim /etc/hosts
127.0.0.1 govanityurls.jianghushinian.cn
```

```go
package main

import (
	"fmt"

	"govanityurls.jianghushinian.cn/blog-go-example/doc/docgo/calculator"
)


func main() {
	fmt.Println(calculator.Add(5, 3))
	fmt.Println(calculator.Subtract(5, 3))
	fmt.Println(calculator.Multiply(5, 3))
	fmt.Println(calculator.Divide(6, 3))
}
```

```shell
$ go get govanityurls.jianghushinian.cn/blog-go-example/doc/docgo/calculator
$ go run main.go                                                            
8
2
15
2
```
