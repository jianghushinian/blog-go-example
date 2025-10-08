package main

import (
	"context"
	"log"
	"os/exec" // 用于执行外部命令（这里用于启动服务器进程）

	"github.com/modelcontextprotocol/go-sdk/mcp" // 导入 MCP Go SDK
)

func main() {
	ctx := context.Background()

	// 创建一个新的 MCP 客户端实例
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// 设置传输层：通过执行外部命令 "myserver" 来启动服务器进程，并使用其 stdin/stdout 进行通信
	transport := &mcp.CommandTransport{Command: exec.Command("./greeting")}

	// 连接到服务器，建立会话（session）
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close() // 确保在函数退出前关闭会话，释放资源

	// 准备工具调用参数
	params := &mcp.CallToolParams{
		Name:      "greet",                        // 工具名称
		Arguments: map[string]any{"name": "江湖十年"}, // 工具参数
	}

	// 调用服务器的工具
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}

	// 处理响应内容
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text) // 类型断言为 TextContent 并获取文本
	}
}
