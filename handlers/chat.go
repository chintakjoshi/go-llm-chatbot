package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chintakjoshi/chintak-chatbot/models"

	"github.com/chintakjoshi/chintak-chatbot/services"

	"github.com/gin-gonic/gin"
)

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

	// Get response from LLM service
	llmReq := services.ChatCompletionRequest{
		Message:   req.Message,
		SessionID: sessionID,
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
		// Log the validation failure
		fmt.Printf("Response validation failed: %s\n", reason)
		fmt.Printf("User message: %s\n", req.Message)
		fmt.Printf("AI response: %s\n", resp.Response)

		// Return a safe, generic response instead of the potentially hallucinated one
		resp.Response = "I'm not sure how to answer that based on my portfolio information. Feel free to ask me about my specific projects, skills, or experience listed on my portfolio website."
	}

	// Simple validation as a fallback
	if !services.SimpleResponseValidation(resp.Response) {
		resp.Response = "I'd be happy to discuss my projects and experience. What would you like to know about my work?"
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
