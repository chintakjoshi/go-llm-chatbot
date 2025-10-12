package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	NvidiaAPIKey       string
	OpenRouterKey      string
	JWTSecret          string
	AllowedOrigins     []string
	RateLimit          int
	RateLimitWindow    int
	NvidiaEndpoint     string
	OpenRouterEndpoint string
	NvidiaModel        string
	OpenRouterModel    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:               mustGetEnv("PORT"),
		NvidiaAPIKey:       mustGetEnv("NVIDIA_API_KEY"),
		OpenRouterKey:      mustGetEnv("OPENROUTER_API_KEY"),
		JWTSecret:          mustGetEnv("JWT_SECRET"),
		AllowedOrigins:     mustGetEnvSlice("ALLOWED_ORIGINS"),
		RateLimit:          getEnvIntWithDefault("RATE_LIMIT", 10),
		RateLimitWindow:    getEnvIntWithDefault("RATE_LIMIT_WINDOW", 60),
		NvidiaEndpoint:     mustGetEnv("NVIDIA_ENDPOINT"),
		OpenRouterEndpoint: mustGetEnv("OPENROUTER_ENDPOINT"),
		NvidiaModel:        mustGetEnv("NVIDIA_MODEL"),
		OpenRouterModel:    mustGetEnv("OPENROUTER_MODEL"),
	}
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func mustGetEnvSlice(key string) []string {
	val := mustGetEnv(key)
	return strings.Split(val, ",")
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s (must be an integer)", key, val)
	}
	return num
}
