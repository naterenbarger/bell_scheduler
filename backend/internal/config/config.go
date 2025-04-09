package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Address     string
	DBPath      string
	JWTSecret   string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	port := getEnvOrDefault("PORT", "8080")
	// Ensure port has a colon prefix if it's just a number
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	cfg := &Config{
		Address:     port,
		DBPath:      getEnvOrDefault("DB_CONNECTION", "bell_scheduler.db"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", ""),
		SMTPHost:    getEnvOrDefault("SMTP_HOST", ""),
		SMTPPort:    getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUser:    getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPass:    getEnvOrDefault("SMTP_PASSWORD", ""),
		SMTPFrom:    getEnvOrDefault("SMTP_FROM", ""),
		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:8080"),
	}

	// Validate required fields
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
