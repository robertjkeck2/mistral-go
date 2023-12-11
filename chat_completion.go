package mistral

import (
	"context"
	"net/http"
)

type FinishReason int

const (
	Stop FinishReason = iota
	Length
)

func (r FinishReason) String() string {
	return [...]string{"stop", "length"}[r]
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeltaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature *float64      `json:"temperature"`
	MaxTokens   *int          `json:"max_tokens"`
	TopP        *float64      `json:"top_p"`
	RandomSeed  *int          `json:"random_seed"`
	Stream      *bool         `json:"stream"`
	SafeMode    *bool         `json:"safe_mode"`
}

type ChatCompletionResponseStreamChoice struct {
	Index        int           `json:"index"`
	Delta        DeltaMessage  `json:"delta"`
	FinishReason *FinishReason `json:"finish_reason"`
}

type ChatCompletionStreamResponse struct {
	ID      string                               `json:"id"`
	Model   string                               `json:"model"`
	Choices []ChatCompletionResponseStreamChoice `json:"choices"`
	Created *int64                               `json:"created"`
	Object  *string                              `json:"object"`
	Usage   *UsageInfo                           `json:"usage"`
}

type ChatCompletionResponseChoice struct {
	Index        int           `json:"index"`
	Message      ChatMessage   `json:"message"`
	FinishReason *FinishReason `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string                         `json:"id"`
	Model   string                         `json:"model"`
	Choices []ChatCompletionResponseChoice `json:"choices"`
	Created int64                          `json:"created"`
	Object  string                         `json:"object"`
	Usage   UsageInfo                      `json:"usage"`
}

func (mc *MistralClient) CreateChatCompletion(ctx context.Context, body ChatCompletionRequest) (resp ChatCompletionResponse, err error) {
	req, err := mc.newRequest(ctx, http.MethodPost, mc.endpoint("/chat/completions"), body)
	if err != nil {
		return
	}

	err = mc.sendRequest(req, &resp)
	return
}

func (mc *MistralClient) CreateChatCompletionStream(ctx context.Context, body ChatCompletionRequest) (resp ChatCompletionStreamResponse, err error) {
	// TODO: implement streaming request and response handling
	return resp, nil
}
