package main

import (
	"context"
	"log"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
	"github.com/go-kratos/blades/memory"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	agent := blades.NewAgent(
		"History Tutor",
		blades.WithModel("deepseek-chat"),
		blades.WithProvider(openai.NewChatProvider(
			option.WithBaseURL("https://api.deepseek.com"),
			option.WithAPIKey("sk-xxx"),
		)),
		blades.WithInstructions("你是一位知识渊博的历史导师。提供详细、准确的历史事件信息。"),
		blades.WithMemory(memory.NewInMemory(10)),
	)
	// Example conversation in memory
	prompt := blades.NewConversation(
		"conversation_123",
		blades.UserMessage("你能告诉我第二次世界大战的起因吗？"),
	)
	result, err := agent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.Text())

	prompt = blades.NewConversation(
		"conversation_123",
		blades.UserMessage("我刚刚问你的问题是什么？"),
	)
	result, err = agent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.Text())
}
