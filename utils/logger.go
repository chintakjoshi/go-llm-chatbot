package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents different log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logFile  *os.File
	logger   *log.Logger
	logLevel LogLevel = INFO
)

// InitLogger initializes the logging system
func InitLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Open log file in append mode, create if it doesn't exist
	file, err := os.OpenFile("logs/backend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFile = file

	// Create logger that writes to both file and stdout
	logger = log.New(file, "", 0) // We'll handle formatting ourselves

	// Set log level from environment variable
	levelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "DEBUG":
		logLevel = DEBUG
	case "WARN":
		logLevel = WARN
	case "ERROR":
		logLevel = ERROR
	case "FATAL":
		logLevel = FATAL
	default:
		logLevel = INFO
	}

	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// getCallerInfo returns the file and line number of the caller
func getCallerInfo() string {
	// Skip 3 levels: getCallerInfo -> log function -> actual caller
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown:0"
	}
	// Get only the filename, not the full path
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// logInternal writes the actual log message
func logInternal(level LogLevel, levelStr string, message string, args ...interface{}) {
	if level < logLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	caller := getCallerInfo()
	formattedMessage := fmt.Sprintf(message, args...)

	logEntry := fmt.Sprintf("[%s] %-5s %s - %s", timestamp, levelStr, caller, formattedMessage)

	// Write to file
	if logger != nil {
		logger.Println(logEntry)
	}

	// Also print to console with colors for different levels
	consoleEntry := fmt.Sprintf("[%s] %s - %s", timestamp, caller, formattedMessage)
	switch level {
	case DEBUG:
		fmt.Printf("\033[36m[DEBUG] %s\033[0m\n", consoleEntry) // Cyan
	case INFO:
		fmt.Printf("\033[32m[INFO]  %s\033[0m\n", consoleEntry) // Green
	case WARN:
		fmt.Printf("\033[33m[WARN]  %s\033[0m\n", consoleEntry) // Yellow
	case ERROR:
		fmt.Printf("\033[31m[ERROR] %s\033[0m\n", consoleEntry) // Red
	case FATAL:
		fmt.Printf("\033[35m[FATAL] %s\033[0m\n", consoleEntry) // Magenta
	}
}

// Log functions for different levels
func Debug(message string, args ...interface{}) {
	logInternal(DEBUG, "DEBUG", message, args...)
}

func Info(message string, args ...interface{}) {
	logInternal(INFO, "INFO", message, args...)
}

func Warn(message string, args ...interface{}) {
	logInternal(WARN, "WARN", message, args...)
}

func Error(message string, args ...interface{}) {
	logInternal(ERROR, "ERROR", message, args...)
}

func Fatal(message string, args ...interface{}) {
	logInternal(FATAL, "FATAL", message, args...)
	os.Exit(1)
}

// Structured logging for specific events
func LogAPIRequest(method, path, clientIP string, statusCode int, duration time.Duration) {
	Info("API Request - %s %s from %s - Status: %d - Duration: %v", method, path, clientIP, statusCode, duration)
}

func LogChatRequest(sessionID, message string, provider string, success bool, duration time.Duration) {
	status := "success"
	if !success {
		status = "failed"
	}
	Info("Chat Request - Session: %s - Provider: %s - Status: %s - Duration: %v - Message: %.50s...",
		sessionID, provider, status, duration, message)
}

func LogProviderSwitch(primaryProvider, fallbackProvider, reason string) {
	Warn("Provider Switch - From: %s, To: %s - Reason: %s", primaryProvider, fallbackProvider, reason)
}

func LogAuthAttempt(clientIP string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	Info("Auth Attempt - Client: %s - Status: %s", clientIP, status)
}
