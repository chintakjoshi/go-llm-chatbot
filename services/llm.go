package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"chintak-chatbot/config"
	"chintak-chatbot/utils"
)

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Message   string
	SessionID string
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	Response  string
	SessionID string
}

// LLMService handles multiple providers with failover
type LLMService struct {
	primaryProvider  Provider
	fallbackProvider Provider
}

// NewLLMService creates a new LLM service with primary and fallback providers
func NewLLMService(cfg *config.Config) *LLMService {
	var primaryProvider, fallbackProvider Provider

	// Initialize NVIDIA as primary if API key is provided
	if cfg.NvidiaAPIKey != "" {
		primaryProvider = NewNvidiaProvider(cfg.NvidiaAPIKey, cfg.NvidiaEndpoint, cfg.NvidiaModel)
		utils.Info("Initialized NVIDIA NIM as primary provider - Model: %s, Endpoint: %s",
			cfg.NvidiaModel, cfg.NvidiaEndpoint)
	} else {
		utils.Warn("NVIDIA API key not provided, skipping NVIDIA provider")
	}

	// Initialize OpenRouter as fallback if API key is provided
	if cfg.OpenRouterKey != "" {
		fallbackProvider = NewOpenRouterProvider(cfg.OpenRouterKey, cfg.OpenRouterEndpoint, cfg.OpenRouterModel)
		utils.Info("Initialized OpenRouter as fallback provider - Model: %s, Endpoint: %s",
			cfg.OpenRouterModel, cfg.OpenRouterEndpoint)
	} else {
		utils.Warn("OpenRouter API key not provided, skipping OpenRouter provider")
	}

	// Validate that at least one provider is configured
	if primaryProvider == nil && fallbackProvider == nil {
		utils.Fatal("No LLM providers configured. Please set at least one API key (NVIDIA_API_KEY or OPENROUTER_API_KEY)")
	}

	return &LLMService{
		primaryProvider:  primaryProvider,
		fallbackProvider: fallbackProvider,
	}
}

// GetChatCompletion gets a chat completion from available providers with failover
func (s *LLMService) GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	startTime := time.Now()

	// Try primary provider first
	if s.primaryProvider != nil {
		utils.Debug("Attempting primary provider: %s", s.primaryProvider.GetName())
		response, err := s.primaryProvider.GetChatCompletion(ctx, req)
		if err == nil {
			duration := time.Since(startTime)
			utils.LogChatRequest(req.SessionID, req.Message, s.primaryProvider.GetName(), true, duration)
			utils.Info("Successfully used %s provider - Duration: %v", s.primaryProvider.GetName(), duration)
			return response, nil
		}

		utils.Warn("Primary provider (%s) failed: %v", s.primaryProvider.GetName(), err)

		// If we have a fallback, try it
		if s.fallbackProvider != nil {
			utils.Info("Attempting failover to %s", s.fallbackProvider.GetName())
			utils.LogProviderSwitch(s.primaryProvider.GetName(), s.fallbackProvider.GetName(), err.Error())

			fallbackResponse, fallbackErr := s.fallbackProvider.GetChatCompletion(ctx, req)
			if fallbackErr == nil {
				duration := time.Since(startTime)
				utils.LogChatRequest(req.SessionID, req.Message, s.fallbackProvider.GetName(), true, duration)
				utils.Info("Successfully used fallback provider (%s) - Duration: %v", s.fallbackProvider.GetName(), duration)
				return fallbackResponse, nil
			}
			utils.Error("Fallback provider (%s) also failed: %v", s.fallbackProvider.GetName(), fallbackErr)

			// Both providers failed
			duration := time.Since(startTime)
			utils.LogChatRequest(req.SessionID, req.Message, "all", false, duration)
			return nil, fmt.Errorf("all providers failed. Primary error: %w, Fallback error: %v", err, fallbackErr)
		}

		// Only primary failed, no fallback
		duration := time.Since(startTime)
		utils.LogChatRequest(req.SessionID, req.Message, s.primaryProvider.GetName(), false, duration)
		return nil, fmt.Errorf("primary provider failed: %w", err)
	}

	// If no primary but we have fallback, use it directly
	if s.fallbackProvider != nil {
		utils.Debug("Using fallback provider as primary: %s", s.fallbackProvider.GetName())
		response, err := s.fallbackProvider.GetChatCompletion(ctx, req)
		if err == nil {
			duration := time.Since(startTime)
			utils.LogChatRequest(req.SessionID, req.Message, s.fallbackProvider.GetName(), true, duration)
			utils.Info("Successfully used fallback provider (%s) as primary - Duration: %v", s.fallbackProvider.GetName(), duration)
			return response, nil
		}
		duration := time.Since(startTime)
		utils.LogChatRequest(req.SessionID, req.Message, s.fallbackProvider.GetName(), false, duration)
		return nil, fmt.Errorf("fallback provider failed: %w", err)
	}

	// This should never happen due to validation in constructor, but just in case
	return nil, fmt.Errorf("no providers available")
}

// ValidateAPIKey checks if at least one provider is working by making a test request
func (s *LLMService) ValidateAPIKey(ctx context.Context) error {
	testReq := ChatCompletionRequest{
		Message: "Hello, please respond with 'OK' if you can read this.",
	}

	response, err := s.GetChatCompletion(ctx, testReq)
	if err != nil {
		return fmt.Errorf("LLM service validation failed: %w", err)
	}

	// Basic validation that we got a response
	if response == nil || response.Response == "" {
		return fmt.Errorf("LLM service returned empty response")
	}

	log.Printf("✅ LLM service validation successful - Providers are working")
	return nil
}

// GetProviderStatus returns the status of available providers
func (s *LLMService) GetProviderStatus() map[string]string {
	status := make(map[string]string)

	if s.primaryProvider != nil {
		status["primary"] = s.primaryProvider.GetName()
	} else {
		status["primary"] = "not configured"
	}

	if s.fallbackProvider != nil {
		status["fallback"] = s.fallbackProvider.GetName()
	} else {
		status["fallback"] = "not configured"
	}

	return status
}

// TestProvider tests a specific provider directly (for debugging)
func (s *LLMService) TestProvider(ctx context.Context, providerName string) error {
	var provider Provider

	switch providerName {
	case "nvidia":
		provider = s.primaryProvider
	case "openrouter":
		provider = s.fallbackProvider
	default:
		return fmt.Errorf("unknown provider: %s", providerName)
	}

	if provider == nil {
		return fmt.Errorf("provider %s is not configured", providerName)
	}

	testReq := ChatCompletionRequest{
		Message: "Test message - please respond with 'OK'",
	}

	_, err := provider.GetChatCompletion(ctx, testReq)
	if err != nil {
		return fmt.Errorf("provider %s test failed: %w", providerName, err)
	}

	return nil
}
