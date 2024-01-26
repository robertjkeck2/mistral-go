package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/robertjkeck2/mistral-go"
)

func main() {
	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		mistral.ChatCompletionRequest{
			Model:    "mistral-tiny",
			Messages: []mistral.ChatMessage{{Role: mistral.RoleUser, Content: "What is the best French cheese?"}},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stream.Close()

	for {
		var response mistral.ChatCompletionStreamResponse
		response, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
