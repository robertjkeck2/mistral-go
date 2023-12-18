package mistral

import (
	"context"
	"net/http"
)

// ChatMessage represents a message in a Mistral chat completion.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	RandomSeed  int           `json:"random_seed,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	SafeMode    bool          `json:"safe_mode,omitempty"`
}

// ChatCompletionResponseChoice represents a choice in a Mistral chat completion.
type ChatCompletionResponseStreamChoice struct {
	Index        int         `json:"index"`
	Delta        ChatMessage `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID      string                               `json:"id"`
	Model   string                               `json:"model"`
	Choices []ChatCompletionResponseStreamChoice `json:"choices"`
	Created *int64                               `json:"created,omitempty"`
	Object  *string                              `json:"object,omitempty"`
	Usage   *UsageInfo                           `json:"usage,omitempty"`
}

// ChatCompletionStream represents a Mistral chat completion stream.
type ChatCompletionStream struct {
	*streamReader[ChatCompletionStreamResponse]
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
	body.Stream = false

	req, err := mc.newRequest(ctx, http.MethodPost, mc.endpoint("/chat/completions"), body)
	if err != nil {
		return
	}

	err = mc.sendRequest(req, &resp)
	return
}

func (mc *MistralClient) CreateChatCompletionStream(ctx context.Context, body ChatCompletionRequest) (stream *ChatCompletionStream, err error) {
	body.Stream = true

	req, err := mc.newRequest(ctx, http.MethodPost, mc.endpoint("/chat/completions"), body)
	if err != nil {
		return
	}

	resp, err := sendRequestStream[ChatCompletionStreamResponse](mc, req)
	if err != nil {
		return
	}

	stream = &ChatCompletionStream{streamReader: resp}
	return
}
