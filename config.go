package mistral

import (
	"net/http"
	"time"
)

const (
	// DefaultMistralURLv1 is the default URL for the Mistral v1 API
	DefaultMistralURLv1 = "https://api.mistral.ai/v1"
)

type ClientConfig struct {
	// ApiKey is the Mistral API key found at https://console.mistral.ai/users/api-keys/
	ApiKey string

	// BaseURL is the base URL for the Mistral API
	BaseURL string

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
		BaseURL:    DefaultMistralURLv1,
		HTTPClient: http.DefaultClient,
		MaxRetries: 5,
		Timeout:    120 * time.Second,
	}
}
