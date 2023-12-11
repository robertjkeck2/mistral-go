package mistral

import "context"

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

func (mc *MistralClient) CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (resp *EmbeddingResponse, err error) {
	return resp, err
}
