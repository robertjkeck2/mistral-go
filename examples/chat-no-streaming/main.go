package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/mistral-go"
)

func main() {
	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		mistral.ChatCompletionRequest{
			Model:    "mistral-tiny",
			Messages: []mistral.ChatMessage{{Role: "user", Content: "What is the best French cheese?"}},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
