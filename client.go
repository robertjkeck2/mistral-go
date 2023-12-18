package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
			bodyReader = reader
		} else {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(bodyBytes)
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
	var (
		retries int
		delay   time.Duration
		err     error
		resp    *http.Response
	)

	retries = mc.config.MaxRetries
	delay = 1 * time.Second
	for retries > 0 {
		resp, err = mc.config.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		if isRetryStatusCode(resp) {
			retries--
			time.Sleep(delay)
			delay *= 2
			continue
		} else if isFailureStatusCode(resp) {
			return mc.handleErrorResp(resp)
		} else {
			break
		}
	}

	defer resp.Body.Close()

	if respBody != nil {
		if err = json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return err
		}
	}

	return nil
}

func (mc *MistralClient) handleErrorResp(resp *http.Response) error {
	var (
		errResp    ErrorResponse
		errMessage ErrorMessage
	)
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(errResp.Message), &errMessage); err != nil {
		return fmt.Errorf("error code %s: %s", errResp.Code, errResp.Message)
	}
	return fmt.Errorf("error code %s: %s", errResp.Code, errMessage.Detail[0].Msg)
}
