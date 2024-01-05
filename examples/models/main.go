package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/mistral-go"
)

func main() {
	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	resp, err := client.ListModels(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(resp.Data) == 0 {
		fmt.Println("No models found")
		return
	}
	for _, model := range resp.Data {
		fmt.Println(model.ID)
	}
}
