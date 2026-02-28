package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Load dedicated public API key for client auth
	publicAPIKey := os.Getenv("PUBLIC_API_KEY")
	if publicAPIKey == "" {
		utils.Fatal("PUBLIC_API_KEY not configured")
	}

	// Initialize services
	llmService := services.NewLLMService(cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authMiddleware, publicAPIKey)
	chatHandler := handlers.NewChatHandler(llmService)

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		utils.Info("Running in release mode")
	} else {
		utils.Info("Running in debug mode")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Trusted proxies
	if err := router.SetTrustedProxies(nil); err != nil {
		utils.Fatal("Failed to set trusted proxies: %v", err)
	}

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit, cfg.RateLimitWindow))
	router.Use(middleware.LoggingMiddleware())

	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/auth", authHandler.Authenticate)

		protected := api.Group("")
		protected.Use(authMiddleware.ValidateToken())
		{
			protected.POST("/chat", chatHandler.Chat)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		status := "ok"
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
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Chintak's Chatbot API is running",
			"version":  "2.0.0",
			"features": "Dual-provider LLM with failover + Structured logging",
		})
	})
	router.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	// Startup logs
	utils.Info("Server starting on port %s", cfg.Port)
	utils.Info("Primary provider: NVIDIA NIM (%s)", cfg.NvidiaModel)
	if cfg.OpenRouterKey != "" {
		utils.Info("Fallback provider: OpenRouter (%s)", cfg.OpenRouterModel)
	}
	utils.Info("Allowed origins: %v", cfg.AllowedOrigins)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown: %v", err)
	}

	utils.Info("Server exited gracefully")
}
