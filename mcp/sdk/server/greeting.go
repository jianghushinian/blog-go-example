package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp" // 导入 MCP Go SDK
)

// Input 结构体定义了工具调用时的输入参数格式
type Input struct {
	// json 标签用于序列化/反序列化，jsonschema 标签提供元数据描述（用于文档或验证）
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

// Output 结构体定义了工具调用后的输出结果格式
type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

// SayHi 是工具的实际处理函数
// 它接收 context、MCP 工具调用请求对象和解析后的 Input 参数。
// 返回三个值：*mcp.CallToolResult（可选的额外结果元数据）、Output（业务输出）和 error（错误）。
func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}

func main() {
	// 创建一个新的 MCP 服务器实例
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	// 向服务器注册一个工具
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

	// 启动服务器，使用标准输入输出（stdio）作为传输层
	// Run 方法会阻塞，直到客户端断开连接或发生错误
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
