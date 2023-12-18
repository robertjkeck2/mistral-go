package mistral

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const MAX_EMPTY_LINES = 300

type streamable interface {
	ChatCompletionStreamResponse
}

type streamReader[T streamable] struct {
	isFinished bool

	reader   *bufio.Reader
	response *http.Response
}

func (stream *streamReader[T]) Recv() (response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	response, err = stream.processLines()
	return
}

func (stream *streamReader[T]) processLines() (T, error) {
	emptyLines := 0

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil {
			_, err := stream.reader.Peek(1)
			if err != nil {
				return *new(T), errors.New("error - message: stream closed unexpectedly")
			}
			return *new(T), readErr
		}

		line := bytes.TrimSpace(rawLine)
		if len(line) == 0 {
			emptyLines++
			if emptyLines > MAX_EMPTY_LINES {
				return *new(T), errors.New("error - message: too many empty lines")
			}
			continue
		}

		if !bytes.HasPrefix(line, []byte("data: ")) {
			return *new(T), errors.New("error - message: invalid line prefix")
		}
		line = line[6:]

		if string(line) == "[DONE]" {
			stream.isFinished = true
			return *new(T), io.EOF
		}

		var response T
		unmarshalErr := json.Unmarshal(line, &response)
		if unmarshalErr != nil {
			return *new(T), unmarshalErr
		}

		return response, nil
	}
}

func (stream *streamReader[T]) Close() {
	stream.response.Body.Close()
}
