package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenRouterProvider struct {
	apiKey   string
	endpoint string
	model    string
	client   *http.Client
}

func NewOpenRouterProvider(apiKey, endpoint, model string) *OpenRouterProvider {
	return &OpenRouterProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (o *OpenRouterProvider) GetName() string {
	return "OpenRouter"
}

func (o *OpenRouterProvider) GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	systemPrompt := GetContextPrompt()
	enhancedMessage := EnhanceUserMessage(req.Message)

	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: enhancedMessage,
		},
	}

	baseReq := BaseRequest{
		Model:       o.model,
		Messages:    messages,
		MaxTokens:   350,
		Temperature: 0.2,
		TopP:        0.7,
		Stream:      false,
	}

	requestBody, err := json.Marshal(baseReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenRouter request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		o.endpoint,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenRouter request: %w", err)
	}

	// Set headers for OpenRouter
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://chintakjoshi.github.io")
	httpReq.Header.Set("X-Title", "Chintak Joshi Portfolio")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("OpenRouter API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var baseResp BaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&baseResp); err != nil {
		return nil, fmt.Errorf("failed to decode OpenRouter response: %w", err)
	}

	if baseResp.Error != nil {
		return nil, fmt.Errorf("OpenRouter API error: %s (type: %s)",
			baseResp.Error.Message, baseResp.Error.Type)
	}

	if len(baseResp.Choices) == 0 || baseResp.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no response content from OpenRouter")
	}

	return &ChatCompletionResponse{
		Response:  baseResp.Choices[0].Message.Content,
		SessionID: req.SessionID,
	}, nil
}
