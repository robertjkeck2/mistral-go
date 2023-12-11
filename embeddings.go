package mistral

type EmbeddingObject struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	ID     string            `json:"id"`
	Model  string            `json:"model"`
	Data   []EmbeddingObject `json:"data"`
	Object string            `json:"object"`
	Usage  UsageInfo         `json:"usage"`
}
