package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// MistralClient is the Mistral API client
type MistralClient struct {
	config ClientConfig
}

// NewMistralClient creates a new Mistral API client
func NewMistralClient(apiKey string) *MistralClient {
	config := DefaultConfig(apiKey)
	return &MistralClient{config: config}
}

func isRetryStatusCode(resp *http.Response) bool {
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func (mc *MistralClient) endpoint(path string) string {
	return mc.config.BaseURL + "/" + mc.config.Version + path
}

func (mc *MistralClient) newRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		if reader, ok := body.(io.Reader); ok {
			body = reader
		} else {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(bodyBytes)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+mc.config.ApiKey)

	return req, nil
}

func (mc *MistralClient) sendRequest(req *http.Request, respBody interface{}) error {
	resp, err := mc.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if isRetryStatusCode(resp) {
		return mc.handleRetryResp(resp)
	}

	if isFailureStatusCode(resp) {
		return mc.handleErrorResp(resp)
	}

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return err
		}
	}

	return nil
}

func (mc *MistralClient) handleRetryResp(resp *http.Response) error {
	return nil
}

func (mc *MistralClient) handleErrorResp(resp *http.Response) error {
	return nil
}
