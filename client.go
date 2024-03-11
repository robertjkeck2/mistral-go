package mistral

import (
	"bufio"
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

// NewMistralClient creates a new Mistral API client with the given config.
func NewMistralClientWithConfig(config ClientConfig) *MistralClient {
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

func sendRequestStream[T streamable](client *MistralClient, req *http.Request) (*streamReader[T], error) {
	var (
		retries int
		delay   time.Duration
		err     error
		resp    *http.Response
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	retries = client.config.MaxRetries
	delay = 1 * time.Second
	for retries > 0 {
		resp, err = client.config.HTTPClient.Do(req)
		if err != nil {
			return new(streamReader[T]), err
		}

		if isRetryStatusCode(resp) {
			retries--
			time.Sleep(delay)
			delay *= 2
			continue
		} else if isFailureStatusCode(resp) {
			return new(streamReader[T]), client.handleErrorResp(resp)
		} else {
			break
		}
	}

	return &streamReader[T]{
		reader:   bufio.NewReader(resp.Body),
		response: resp,
	}, nil
}

func (mc *MistralClient) handleErrorResp(resp *http.Response) error {
	var (
		errString  string
		errResp    ErrorResponse
		errMessage ErrorMessage
	)

	errorBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errString = "error - message: failed to read error response body"
	}
	errString = string(errorBytes)

	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("error - message: %s", errString)
	}

	if err := json.Unmarshal([]byte(errResp.Message), &errMessage); err != nil {
		return fmt.Errorf("error - code: %s, message: %s", errResp.Code, errResp.Message)
	}

	return fmt.Errorf("error - code: %s, message: %s", errResp.Code, errMessage.Detail[0].Msg)
}
