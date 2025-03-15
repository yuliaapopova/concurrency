package replication

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
)

type Request struct {
	LastSegmentName string
}

func NewRequest(lastSegmentName string) *Request {
	return &Request{
		LastSegmentName: lastSegmentName,
	}
}

type Response struct {
	Success     bool
	SegmentName string
	SegmentData []byte
}

func NewResponse(success bool, segmentName string, segmentData []byte) *Response {
	return &Response{
		Success:     success,
		SegmentName: segmentName,
		SegmentData: segmentData,
	}
}

func EncodeRequest(request *Request) ([]byte, error) {
	if request == nil {
		return nil, errors.New("encodeRequest: request is nil")
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	return buf.Bytes(), nil
}

func DecodeRequest(ctx context.Context, request *Request, data []byte) error {
	if request == nil {
		return errors.New("decodeRequest: request is nil")
	}

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(request); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}

	return nil
}

func EncodeResponse(ctx context.Context, response *Response) ([]byte, error) {
	if response == nil {
		return nil, errors.New("encodeResponse: response is nil")
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(response); err != nil {
		return nil, fmt.Errorf("failed to encode response: %w", err)
	}

	return buf.Bytes(), nil
}

func DecodeResponse(response *Response, data []byte) error {
	if response == nil {
		return errors.New("decodeResponse: response is nil")
	}

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
