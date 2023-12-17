package mistral

import (
	"context"
	"net/http"
)

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
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float64      `json:"top_p,omitempty"`
	RandomSeed  *int          `json:"random_seed,omitempty"`
	Stream      *bool         `json:"stream,omitempty"`
	SafeMode    *bool         `json:"safe_mode,omitempty"`
}

type ChatCompletionResponseStreamChoice struct {
	Index        int          `json:"index"`
	Delta        DeltaMessage `json:"delta"`
	FinishReason *string      `json:"finish_reason,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID      string                               `json:"id"`
	Model   string                               `json:"model"`
	Choices []ChatCompletionResponseStreamChoice `json:"choices"`
	Created *int64                               `json:"created,omitempty"`
	Object  *string                              `json:"object,omitempty"`
	Usage   *UsageInfo                           `json:"usage,omitempty"`
}

type ChatCompletionResponseChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason *string     `json:"finish_reason,omitempty"`
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
