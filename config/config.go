package config

import (
	"log"
	"net/url"
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

	cfg := &Config{
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

	validateConfig(cfg)
	return cfg
}

func validateConfig(cfg *Config) {
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	if len(cfg.AllowedOrigins) == 0 {
		log.Fatal("ALLOWED_ORIGINS must not be empty")
	}

	for k, v := range map[string]string{
		"NVIDIA_ENDPOINT":     cfg.NvidiaEndpoint,
		"OPENROUTER_ENDPOINT": cfg.OpenRouterEndpoint,
	} {
		if _, err := url.ParseRequestURI(v); err != nil {
			log.Fatalf("Invalid %s: %v", k, err)
		}
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
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
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
