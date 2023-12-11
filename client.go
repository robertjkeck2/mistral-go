package mistral

// Client is the Mistral API client
type Client struct {
	config ClientConfig
}

// NewClient creates a new Mistral API client
func NewClient(apiKey string) *Client {
	config := DefaultConfig(apiKey)
	return &Client{config: config}
}
