package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// OpenRouterRequest represents the request structure for OpenRouter API
type OpenRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []OpenRouterMessage `json:"messages"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
	TopP        float64             `json:"top_p,omitempty"`
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response structure from OpenRouter API
type OpenRouterResponse struct {
	Choices []OpenRouterChoice `json:"choices"`
	Error   *OpenRouterError   `json:"error,omitempty"`
}

type OpenRouterChoice struct {
	Message OpenRouterMessage `json:"message"`
}

type OpenRouterError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type LLMService struct {
	apiKey     string
	httpClient *http.Client
}

func NewLLMService(apiKey string) *LLMService {
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY is required")
	}

	return &LLMService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ChatCompletionRequest struct {
	Message   string
	SessionID string
}

type ChatCompletionResponse struct {
	Response  string
	SessionID string
}

func (s *LLMService) GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	systemPrompt := GetContextPrompt()
	enhancedMessage := EnhanceUserMessage(req.Message)

	// Prepare messages for OpenRouter
	messages := []OpenRouterMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: enhancedMessage,
		},
	}

	openRouterReq := OpenRouterRequest{
		Model:       "z-ai/glm-4.5-air:free", // Changed model
		Messages:    messages,
		MaxTokens:   5000, // Further reduced to prevent long responses
		Temperature: 0.2,  // Much lower for more deterministic responses
		TopP:        0.7,  // Adjusted for better response quality
	}

	requestBody, err := json.Marshal(openRouterReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for OpenRouter
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	// OpenRouter requires these headers
	httpReq.Header.Set("HTTP-Referer", "https://chintakjoshi.github.io")
	httpReq.Header.Set("X-Title", "Chintak Joshi Portfolio")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("OpenRouter API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter API returned status: %d", resp.StatusCode)
	}

	var openRouterResp OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if openRouterResp.Error != nil {
		return nil, fmt.Errorf("OpenRouter API error: %s (type: %s)",
			openRouterResp.Error.Message, openRouterResp.Error.Type)
	}

	if len(openRouterResp.Choices) == 0 || openRouterResp.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no response content from OpenRouter")
	}

	return &ChatCompletionResponse{
		Response:  openRouterResp.Choices[0].Message.Content,
		SessionID: req.SessionID,
	}, nil
}

// ValidateAPIKey checks if the OpenRouter API key is valid
func (s *LLMService) ValidateAPIKey(ctx context.Context) error {
	// For OpenRouter, we'll do a simple test call
	testReq := ChatCompletionRequest{
		Message: "Hello",
	}

	_, err := s.GetChatCompletion(ctx, testReq)
	if err != nil {
		return fmt.Errorf("OpenRouter API key validation failed: %w", err)
	}
	return nil
}
