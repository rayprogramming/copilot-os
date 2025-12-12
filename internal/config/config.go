package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds configuration for the MCP server and orchestrator.
type Config struct {
	// RepoRoot is the path to the repository containing agents.
	RepoRoot string

	// LogLevel is the logging level (debug, info, warn, error).
	LogLevel string

	// CacheEnabled controls whether result caching is enabled.
	CacheEnabled bool

	// CLITimeout is the timeout for Copilot CLI calls.
	CLITimeout time.Duration
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() *Config {
	cfg := &Config{
		RepoRoot:     getEnv("REPO_ROOT", "."),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		CacheEnabled: getEnvBool("CACHE_ENABLED", true),
		CLITimeout:   getEnvDuration("COPILOT_CLI_TIMEOUT", 300*time.Second),
	}
	return cfg
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvBool retrieves a boolean environment variable or returns a default value.
func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

// getEnvDuration retrieves a duration environment variable or returns a default value.
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
