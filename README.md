# Mistral Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/robertjkeck2/mistral-go.svg)](https://pkg.go.dev/github.com/robertjkeck2/mistral-go)
[![Go Report Card](https://goreportcard.com/report/github.com/robertjkeck2/go-mistral)](https://goreportcard.com/report/github.com/robertjkeck2/go-mistral)

Unofficial Golang client library for Mistral AI platform

## Installation

`go get github.com/robertjkeck2/mistral-go`

_Requires Go 1.21 or later_

## Usage

### Mistral Chat Completion example:

```
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/go-mistral"
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
```

### Getting a Mistral API key:

1. Sign up for a Mistral account at https://console.mistral.ai/.
2. If you don't already have one, create an account.
3. In the `Your subscription` section, click `Manage`.
4. Update your billing information in the `Billing information` section.
5. Create a new subscription in the `Current subscription` section.
6. Go to the `API Keys` tab in the left sidebar.
7. Click `Generate a new key` and copy the key. Save the key somewhere safe and do not share this key with anyone.

### Other examples:

<details>
<summary>Chat Completion Streaming</summary>

```
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/robertjkeck2/go-mistral"
)

func main() {
	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	stream, err := client.CreateChatCompletionStream(
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
```

</details>

<details>
<summary>Embeddings</summary>

```
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/go-mistral"
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
```

</details>

<details>
<summary>List Models</summary>

```
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/robertjkeck2/go-mistral"
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
```

</details>

### Mistral API documentation

https://docs.mistral.ai/api
