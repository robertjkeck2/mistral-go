package mistral

// MistralClient is the Mistral API client
type MistralClient struct {
	config ClientConfig
}

// NewMistralClient creates a new Mistral API client
func NewMistralClient(apiKey string) *MistralClient {
	config := DefaultConfig(apiKey)
	return &MistralClient{config: config}
}
