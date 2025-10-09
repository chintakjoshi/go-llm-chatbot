package handlers

import (
	"context"
	"net/http"
	"time"

	"chintak-chatbot/models"
	"chintak-chatbot/services"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	llmService *services.LLMService // Changed from openAIService
}

func NewChatHandler(llmService *services.LLMService) *ChatHandler { // Changed parameter type
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

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Message is required",
		})
		return
	}

	// Generate session ID if not provided
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Get response from OpenAI
	openAIReq := services.ChatCompletionRequest{
		Message:   req.Message,
		SessionID: sessionID,
	}

	resp, err := h.llmService.GetChatCompletion(ctx, openAIReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get response",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ChatResponse{
		Response:  resp.Response,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
	})
}

func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
