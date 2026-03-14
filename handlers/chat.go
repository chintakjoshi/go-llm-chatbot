package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chintakjoshi/chintak-chatbot/models"
	"github.com/chintakjoshi/chintak-chatbot/services"
	"github.com/gin-gonic/gin"
)

const maxMessageLength = 2000 // prevent prompt overflow / abuse

type ChatHandler struct {
	llmService *services.LLMService
}

func NewChatHandler(llmService *services.LLMService) *ChatHandler {
	return &ChatHandler{
		llmService: llmService,
	}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Message is required",
		})
		return
	}
	if len(req.Message) > maxMessageLength {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf("Message too long (max %d characters)", maxMessageLength),
		})
		return
	}

	// Generate session ID if not provided (cryptographically secure)
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	decision := services.ApplyGuardrailsWithSession(sessionID, req.Message)
	if decision.DirectResponse != "" {
		c.JSON(http.StatusOK, models.ChatResponse{
			Response:  decision.DirectResponse,
			SessionID: sessionID,
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Get response from LLM service
	llmReq := services.ChatCompletionRequest{
		Message:      req.Message,
		SessionID:    sessionID,
		SystemPrompt: decision.SystemPrompt,
		UserPrompt:   decision.UserPrompt,
	}

	resp, err := h.llmService.GetChatCompletion(ctx, llmReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get response",
			Message: err.Error(),
		})
		return
	}

	// Validate the response
	if isValid, reason := services.ValidateResponse(resp.Response, req.Message); !isValid {
		fmt.Printf("Response validation failed: %s\n", reason)
		fmt.Printf("User message: %s\n", req.Message)
		fmt.Printf("AI response: %s\n", resp.Response)

		resp.Response = services.PortfolioOnlyResponse
	}

	// Simple validation as a fallback
	if !services.SimpleResponseValidation(resp.Response) {
		resp.Response = services.PortfolioOnlyResponse
	}

	c.JSON(http.StatusOK, models.ChatResponse{
		Response:  resp.Response,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
	})
}

func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// If crypto/rand fails, fall back to time-based uniqueness
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(b)
}
