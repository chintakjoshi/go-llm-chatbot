package services

import (
	"context"
)

// Provider defines the interface for LLM providers
type Provider interface {
	GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
	GetName() string
}

// BaseRequest is the common request structure for both providers
type BaseRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BaseResponse is the common response structure
type BaseResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
}
