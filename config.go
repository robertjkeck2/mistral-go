package mistral

import (
	"net/http"
	"time"
)

const (
	// DefaultMistralURL is the default URL for the Mistral API
	DefaultMistralURL = "https://api.mistral.ai"
)

type ClientConfig struct {
	// ApiKey is the Mistral API key found at https://console.mistral.ai/users/api-keys/
	ApiKey string

	// BaseURL is the base URL for the Mistral API
	BaseURL string

	// Version is the version of the Mistral API
	Version string

	// HTTPClient is the HTTP client used for requests
	HTTPClient *http.Client

	// MaxRetries is the maximum number of retries for a request
	MaxRetries int

	// Timeout is the timeout for a request
	Timeout time.Duration
}

func DefaultConfig(apiKey string) ClientConfig {
	return ClientConfig{
		ApiKey:     apiKey,
		BaseURL:    DefaultMistralURL,
		Version:    "v1",
		HTTPClient: http.DefaultClient,
		MaxRetries: 5,
		Timeout:    120 * time.Second,
	}
}
