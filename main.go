package main

import (
	"log"
	"os"
	"time"

	"chintak-chatbot/config"
	"chintak-chatbot/handlers"
	"chintak-chatbot/middleware"
	"chintak-chatbot/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.OpenRouterKey == "" {
		log.Fatal("OPENROUTER_API_KEY is required")
	}

	// Initialize services - changed from OpenAI to LLM service
	llmService := services.NewLLMService(cfg.OpenRouterKey) // Changed
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authMiddleware, cfg.OpenRouterKey)
	chatHandler := handlers.NewChatHandler(llmService) // This now uses LLMService

	// Create Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	router.Use(middleware.RateLimitMiddleware())

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

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "chintak-chatbot",
			"time":    time.Now().Unix(),
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Chintak's Chatbot API is running",
			"version": "1.0.0",
			"model":   "deepseek/deepseek-chat-v3.1:free",
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Using OpenRouter with DeepSeek model")
	log.Printf("Allowed origins: %v", cfg.AllowedOrigins)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
