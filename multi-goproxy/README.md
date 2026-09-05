# GOPROXY 逗号 vs 管道 测试

这个项目用于验证 `GOPROXY` 配置中逗号 (`,`) 和管道 (`|`) 分隔符的区别。

## 核心区别

| 分隔符 | 回退条件 | 说明 |
|--------|----------|------|
| `,` (逗号) | 只在 404/410 时回退 | 遇到其他错误（500/超时）停止 |
| `\|` (管道) | 任何错误都回退 | 遇到任何错误都尝试下一个代理 |

## 测试环境

- Go 1.21+
- macOS / Linux（依赖 bash、`/dev/tcp`）
- 只访问 `127.0.0.1`,不依赖外网

## 快速开始

### 1. 运行自动测试（推荐）

```bash
cd multi-goproxy
chmod +x test.sh
./test.sh
```

脚本启动 3 个本地代理、跑 4 个用例并做交叉断言,结束时自动清理临时目录与代理进程:

```
用例: 逗号模式: 404 后遇到 500 应停止
  [通过] go mod download 退出码=1(符合预期)
用例: 管道模式: 404 后遇到 500 应继续尝试
  [通过] go mod download 退出码=0(符合预期)
...
结果: 通过 4 / 4
```

### 2. 手动验证

需要 3 个终端,分别启动 3 个代理;再开 1 个终端跑 `go mod download -x`。

#### 启动 3 个代理服务器

```bash
# 终端 1: 代理 p1 —— 返回 404(该代理上找不到模块)
PROXY_NAME=p1 PROXY_PORT=18081 PROXY_STATUS=404 ./proxy-server

# 终端 2: 代理 p2 —— 返回 500(代理服务端出错)
PROXY_NAME=p2 PROXY_PORT=18082 PROXY_STATUS=500 ./proxy-server

# 终端 3: 代理 p3 —— 返回 200(健康代理,按 GOPROXY 协议真实提供模块)
PROXY_NAME=p3 PROXY_PORT=18083 PROXY_STATUS=200 ./proxy-server
```

> 200 代理会根据请求的模块路径动态返回合法的 `list/info/mod/zip` 内容,
> 让 `go mod download` 能真正下载成功;`example.com` 这类父路径探测也会被正常处理。

#### 观察回退过程

选 `go mod download`(钉死版本)而不是 `go get`,是因为钉死版本后 go 只取
`@v/<version>.{info,mod,zip}` 三个端点,请求序列规整可预期;`go get` 还会做版本
解析(`@latest`/`list`)与父模块探测,请求更杂,并会改写 `go.mod`。

带 `-x` 参数时,go 会把**每一次代理请求及其 HTTP 结果**逐条打到 stderr,
回退链(先试谁、为什么回退、最后停在哪)直接可见。

先隔离模块缓存(命中缓存就不会访问代理),成功过一次后想重跑要**换一个新的缓存目录**:

```bash
export GOSUMDB=off GOINSECURE='*' GOTOOLCHAIN=local
export GOMODCACHE="$(mktemp -d)"   # 每次验证前换成新目录,保证真的去访问代理
```

**测试 1: 逗号模式（p1=404 → p2=500 应停止）**

```bash
export GOPROXY="http://127.0.0.1:18081,http://127.0.0.1:18082,http://127.0.0.1:18083"
cd testmodule && go mod download -x example.com/dep@v1.0.0
```

预期:非 0 退出。`-x` 显示 p1 404 → 回退 p2 → p2 500 → **停止**,全程没试 p3:

```
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.info
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.info: 404 Not Found
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.info
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.info: 500 Internal Server Error
go: example.com/dep@v1.0.0: reading http://127.0.0.1:18082/...@v/v1.0.0.info: 500 Internal Server Error
```

**测试 2: 管道模式（p1=404 → p2=500 → p3=200 应成功）**

```bash
export GOPROXY="http://127.0.0.1:18081|http://127.0.0.1:18082|http://127.0.0.1:18083"
cd testmodule && go mod download -x example.com/dep@v1.0.0
```

预期:退出码 0。`.info/.mod/.zip` 每个文件都从 p1 一路试到 p3:

```
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.info
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.info: 404 Not Found (0.004s)
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.info
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.info: 500 Internal Server Error (0.001s)
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.info
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.info: 200 OK (0.001s)
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.mod
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.mod: 404 Not Found (0.000s)
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.mod
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.mod: 500 Internal Server Error (0.000s)
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.mod
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.mod: 200 OK (0.000s)
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.zip
# get http://127.0.0.1:18081/example.com/dep/@v/v1.0.0.zip: 404 Not Found (0.000s)
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.zip
# get http://127.0.0.1:18082/example.com/dep/@v/v1.0.0.zip: 500 Internal Server Error (0.000s)
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.zip
# get http://127.0.0.1:18083/example.com/dep/@v/v1.0.0.zip: 200 OK (0.000s)
```

把首位置换成 p2(500)再做一遍,得到:
- 逗号 `http://127.0.0.1:18082,...`:失败——首位 500 立即停止;
- 管道 `http://127.0.0.1:18082|...`:成功——500 后继续 p1(404)再 p3(200)。

#### 手动清理

3 个代理在其终端里按 `Ctrl-C` 即可。若残留进程占用端口:

```bash
lsof -nP -iTCP:18081 -iTCP:18082 -iTCP:18083 -sTCP:LISTEN
kill <上述 PID>
```

## 预期结果

| 测试 | GOPROXY 配置 | p1 | p2 | p3 | 预期结果 |
|------|--------------|----|----|----|----------|
| 1 | `p1,p2,p3` | 404 | 500 | 200 | 失败（逗号模式在 500 停止） |
| 2 | `p1\|p2\|p3` | 404 | 500 | 200 | 成功（管道模式继续尝试） |
| 3 | `p2,p1,p3` | 500 | 404 | 200 | 失败（逗号模式在 500 停止） |
| 4 | `p2\|p1\|p3` | 500 | 404 | 200 | 成功（管道模式继续尝试） |

## 目录与脚本说明

- `main.go` —— 模拟 Go module proxy,用 `PROXY_STATUS` 控制行为(见下节);
- `test.sh` —— 自动化验证脚本:编译代理 → 端口预检 → 4 个用例断言 → 汇总清理;
- `testmodule/` —— 手动验证用的示例模块,`require example.com/dep v1.0.0`。

`test.sh` 的每个用例都在独立临时目录里隔离 `GOMODCACHE/GOPATH/GOCACHE/GOENV`,
并做两项交叉断言:1) `go mod download` 退出码是否符合预期;2) 健康代理 p3(200)是否收到请求。
仅当两者都符合时用例才判为通过,避免"退出码碰巧正确但根本没走代理"的假阳性。

## 代理服务器参数

`main.go` 支持命令行参数或环境变量(不带参数时读环境变量):

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `PROXY_NAME` | 代理名称(用于日志) | `proxy` |
| `PROXY_PORT` | 监听端口 | `8080` |
| `PROXY_STATUS` | 返回的 HTTP 状态码 | `200` |

| 命令行参数 | 说明 | 默认值 |
|------------|------|--------|
| `-name` | 代理名称 | `proxy` |
| `-port` | 监听端口 | `8080` |
| `-status` | HTTP 状态码 | `200` |

### 支持的状态码

- `200`: 健康代理,按 GOPROXY 协议动态提供模块内容(list/info/mod/zip)
- `404`: Not Found(触发逗号模式回退)
- `500`/`502`/`503`: 服务端错误(逗号模式停止)

## 实际应用场景

### 场景 1: 公司代理 + 公共代理

```bash
# 公司代理可能返回 500,使用管道模式确保回退到公共代理
export GOPROXY="https://corp-proxy.com|https://goproxy.cn|direct"
```

### 场景 2: 国内代理 + 国际代理

```bash
# 国内代理优先,404 时回退到国际代理
export GOPROXY="https://goproxy.cn,https://proxy.golang.org,direct"
```

## 参考

- [Go Modules Reference - GOPROXY](https://go.dev/ref/mod#goproxy-protocol)
- [goproxy.cn](https://goproxy.cn)
- [proxy.golang.org](https://proxy.golang.org)
