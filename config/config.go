package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	OpenRouterKey   string // Changed from OpenAIKey
	JWTSecret       string
	AllowedOrigins  []string
	RateLimit       int
	RateLimitWindow int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:            getEnv("PORT", "8080"),
		OpenRouterKey:   getEnv("OPENROUTER_API_KEY", ""), // Changed
		JWTSecret:       getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		AllowedOrigins:  getEnvSlice("ALLOWED_ORIGINS", []string{"https://chintakjoshi.github.io", "http://localhost:3000"}),
		RateLimit:       getEnvInt("RATE_LIMIT", 10),
		RateLimitWindow: getEnvInt("RATE_LIMIT_WINDOW", 60),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
