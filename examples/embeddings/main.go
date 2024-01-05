package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/mistral-go"
)

func main() {
	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	resp, err := client.CreateEmbedding(
		context.Background(),
		mistral.EmbeddingRequest{
			Model: "mistral-embed",
			Input: []string{"What is the best French cheese?"},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Data[0].Embedding)
}
