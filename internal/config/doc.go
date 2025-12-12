// Package config provides configuration management for the CopilotOS server.
//
// This package handles loading and managing configuration from environment variables,
// with sensible defaults for all settings.
//
// # Configuration Sources
//
// All configuration is loaded from environment variables. If an environment variable
// is not set, a default value is used. This approach follows the twelve-factor app
// methodology for configuration management.
//
// # Configuration Options
//
// The following environment variables are supported:
//
//	REPO_ROOT           - Path to the repository containing agents (default: ".")
//	LOG_LEVEL           - Logging level: debug, info, warn, error (default: "info")
//	CACHE_ENABLED       - Enable result caching: true, false (default: true)
//	COPILOT_CLI_TIMEOUT - Timeout for Copilot CLI calls (default: 300s)
//
// Usage Example
//
//	// Load configuration from environment
//	cfg := config.LoadFromEnv()
//
//	// Access configuration values
//	fmt.Printf("Repository root: %s\n", cfg.RepoRoot)
//	fmt.Printf("Log level: %s\n", cfg.LogLevel)
//	fmt.Printf("Cache enabled: %t\n", cfg.CacheEnabled)
//	fmt.Printf("CLI timeout: %s\n", cfg.CLITimeout)
//
// # Setting Environment Variables
//
// You can set environment variables in several ways:
//
// 1. Shell export:
//
//	export REPO_ROOT=/path/to/repo
//	export LOG_LEVEL=debug
//	export CACHE_ENABLED=false
//	export COPILOT_CLI_TIMEOUT=5m
//
// 2. .env file (if using a tool like godotenv):
//
//	REPO_ROOT=/path/to/repo
//	LOG_LEVEL=debug
//	CACHE_ENABLED=false
//	COPILOT_CLI_TIMEOUT=5m
//
// 3. Docker/Kubernetes environment:
//
//	docker run -e REPO_ROOT=/repo -e LOG_LEVEL=debug ...
//
// # Default Values
//
// The default values are chosen for development and local testing:
//   - REPO_ROOT: "." (current directory)
//   - LOG_LEVEL: "info" (balanced logging)
//   - CACHE_ENABLED: true (improve performance)
//   - COPILOT_CLI_TIMEOUT: 300s (5 minutes, accommodates slow operations)
//
// For production deployments, consider adjusting:
//   - LOG_LEVEL: "warn" or "error" (reduce log volume)
//   - COPILOT_CLI_TIMEOUT: Lower value if agents are expected to be fast
//
// # Type Safety
//
// The package provides type-safe parsing for environment variables:
//   - Strings: Direct string values
//   - Booleans: Parsed with strconv.ParseBool (accepts: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)
//   - Durations: Parsed with time.ParseDuration (e.g., "5m", "30s", "1h30m")
//
// If parsing fails, the default value is returned silently. This ensures the
// server can always start with valid configuration.
package config
