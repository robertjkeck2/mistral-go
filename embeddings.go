package mistral

import (
	"context"
	"net/http"
)

type EmbeddingObject struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingResponse struct {
	ID     string            `json:"id"`
	Model  string            `json:"model"`
	Data   []EmbeddingObject `json:"data"`
	Object string            `json:"object"`
	Usage  UsageInfo         `json:"usage"`
}

func (mc *MistralClient) CreateEmbedding(ctx context.Context, body *EmbeddingRequest) (resp *EmbeddingResponse, err error) {
	req, err := mc.newRequest(ctx, http.MethodPost, mc.endpoint("/embeddings"), body)
	if err != nil {
		return
	}

	err = mc.sendRequest(req, &resp)
	return
}
