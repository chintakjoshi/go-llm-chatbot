package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chintakjoshi/chintak-chatbot/config"
	"github.com/chintakjoshi/chintak-chatbot/handlers"
	"github.com/chintakjoshi/chintak-chatbot/middleware"
	"github.com/chintakjoshi/chintak-chatbot/services"
	"github.com/chintakjoshi/chintak-chatbot/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logging
	if err := utils.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer utils.CloseLogger()

	utils.Info("Starting Chintak Chatbot Backend...")

	// Load configuration
	cfg := config.Load()
	utils.Info("Configuration loaded successfully")

	// Validate required configuration
	if cfg.NvidiaAPIKey == "" && cfg.OpenRouterKey == "" {
		utils.Fatal("No LLM providers configured. Please set at least one API key (NVIDIA_API_KEY or OPENROUTER_API_KEY)")
	}

	// Initialize services with the new LLM service
	utils.Info("Initializing LLM service...")
	llmService := services.NewLLMService(cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authMiddleware, cfg.NvidiaAPIKey)
	chatHandler := handlers.NewChatHandler(llmService)

	// Create Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		utils.Info("Running in release mode")
	} else {
		utils.Info("Running in debug mode")
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.LoggingMiddleware()) // Add logging middleware

	// Routes
	api := router.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth", authHandler.Authenticate)
		api.GET("/auth/simple", authHandler.SimpleAuth)

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.ValidateToken())
		{
			protected.POST("/chat", chatHandler.Chat)
		}
	}

	// Health check with provider status
	router.GET("/health", func(c *gin.Context) {
		status := "ok"

		// Quick provider test
		testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := llmService.ValidateAPIKey(testCtx); err != nil {
			status = "degraded"
			utils.Warn("Health check warning: %v", err)
		}

		c.JSON(200, gin.H{
			"status":    status,
			"service":   "chintak-chatbot",
			"time":      time.Now().Unix(),
			"providers": "NVIDIA NIM (primary) + OpenRouter (fallback)",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Chintak's Chatbot API is running",
			"version":  "2.0.0",
			"features": "Dual-provider LLM with failover + Structured logging",
		})
	})

	utils.Info("Server starting on port %s", cfg.Port)
	utils.Info("Primary provider: NVIDIA NIM (%s)", cfg.NvidiaModel)
	if cfg.OpenRouterKey != "" {
		utils.Info("Fallback provider: OpenRouter (%s)", cfg.OpenRouterModel)
	}
	utils.Info("Allowed origins: %v", cfg.AllowedOrigins)

	if err := router.Run(":" + cfg.Port); err != nil {
		utils.Fatal("Failed to start server: %v", err)
	}
}
