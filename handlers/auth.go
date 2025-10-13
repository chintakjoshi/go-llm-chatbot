package handlers

import (
	"net/http"

	"github.com/chintakjoshi/chintak-chatbot/middleware"
	"github.com/chintakjoshi/chintak-chatbot/models"
	"github.com/chintakjoshi/chintak-chatbot/utils"
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

	// The server must be configured with a validation key; no fallback in prod.
	if h.validAPIKey == "" {
		utils.Error("Authentication attempted but VALID_API_KEY is not configured")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Server misconfiguration",
		})
		return
	}

	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("Invalid auth request from %s: %v", clientIP, err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Constant-time comparison would be ideal if available; here we do a simple check.
	if req.APIKey != h.validAPIKey {
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
