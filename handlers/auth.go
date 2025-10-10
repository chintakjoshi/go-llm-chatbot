package handlers

import (
	"net/http"

	"chintak-chatbot/middleware"
	"chintak-chatbot/models"
	"chintak-chatbot/utils"

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
	clientIP := c.ClientIP()

	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("Invalid auth request from %s: %v", clientIP, err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// For now, we'll use a simple API key validation
	expectedAPIKey := h.validAPIKey
	if expectedAPIKey == "" {
		// Fallback for development
		expectedAPIKey = "portfolio-chatbot-key"
	}

	if req.APIKey != expectedAPIKey {
		utils.Warn("Failed authentication attempt from %s - Invalid API key", clientIP)
		utils.LogAuthAttempt(clientIP, false)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid API key",
		})
		return
	}

	// Generate JWT token
	token, err := h.authMiddleware.GenerateToken(req.APIKey)
	if err != nil {
		utils.Error("Failed to generate token for %s: %v", clientIP, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	utils.Info("Successful authentication from %s", clientIP)
	utils.LogAuthAttempt(clientIP, true)

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 hours in seconds
	})
}

func (h *AuthHandler) SimpleAuth(c *gin.Context) {
	clientIP := c.ClientIP()

	apiKey := "portfolio-chatbot-key" // You can change this

	token, err := h.authMiddleware.GenerateToken(apiKey)
	if err != nil {
		utils.Error("Failed to generate simple auth token for %s: %v", clientIP, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	utils.Info("Successful simple authentication from %s", clientIP)
	utils.LogAuthAttempt(clientIP, true)

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400,
	})
}
