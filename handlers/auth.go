package handlers

import (
	"net/http"

	"chintak-chatbot/middleware"
	"chintak-chatbot/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authMiddleware *middleware.AuthMiddleware
	validAPIKey    string
}

func NewAuthHandler(authMiddleware *middleware.AuthMiddleware, validAPIKey string) *AuthHandler {
	return &AuthHandler{
		authMiddleware: authMiddleware,
		validAPIKey:    validAPIKey,
	}
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// For now, we'll use a simple API key validation
	// In production, you might want a more secure approach
	expectedAPIKey := h.validAPIKey
	if expectedAPIKey == "" {
		// Fallback for development
		expectedAPIKey = "portfolio-chatbot-key"
	}

	if req.APIKey != expectedAPIKey {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid API key",
		})
		return
	}

	// Generate JWT token
	token, err := h.authMiddleware.GenerateToken(req.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 hours in seconds
	})
}

// SimpleAuth is a simpler version that doesn't require request body
func (h *AuthHandler) SimpleAuth(c *gin.Context) {
	// For GitHub Pages, we might want a simpler approach
	// This generates a token without requiring an API key in request body

	apiKey := "portfolio-chatbot-key" // You can change this

	token, err := h.authMiddleware.GenerateToken(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400,
	})
}
