package main

import (
	"context"
	"log"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	agent := blades.NewAgent(
		"Template Agent",
		blades.WithModel("deepseek-chat"),
		blades.WithProvider(openai.NewChatProvider(
			option.WithBaseURL("https://api.deepseek.com"),
			option.WithAPIKey("sk-xxx"),
		)),
	)

	// Define templates and params
	params := map[string]any{
		"topic":    "人工智能的未来",
		"audience": "普通读者",
	}

	// Build prompt using the template builder
	// Note: Use exported methods when calling from another package.
	prompt, err := blades.NewPromptTemplate().
		System("请将 {{.topic}} 总结为三个关键点。", params).
		User("对 {{.audience}} 受众做出简洁准确的回复。", params).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generated Prompt:", prompt.String())

	// Run the agent with the templated prompt
	stream, err := agent.RunStream(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	for stream.Next() {
		chunk, err := stream.Current()
		if err != nil {
			log.Fatalf("stream recv error: %v", err)
		}
		log.Print(chunk.Text())
	}
}
