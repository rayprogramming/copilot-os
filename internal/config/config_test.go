package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPO_ROOT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("CACHE_ENABLED")
	os.Unsetenv("COPILOT_CLI_TIMEOUT")

	cfg := LoadFromEnv()

	if cfg.RepoRoot != "." {
		t.Errorf("expected default RepoRoot '.', got %q", cfg.RepoRoot)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected default LogLevel 'info', got %q", cfg.LogLevel)
	}

	if cfg.CacheEnabled != true {
		t.Errorf("expected default CacheEnabled true, got %v", cfg.CacheEnabled)
	}

	if cfg.CLITimeout != 300*time.Second {
		t.Errorf("expected default CLITimeout 300s, got %v", cfg.CLITimeout)
	}
}

func TestLoadFromEnv_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("REPO_ROOT", "/custom/path")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("CACHE_ENABLED", "false")
	os.Setenv("COPILOT_CLI_TIMEOUT", "60s")
	defer func() {
		os.Unsetenv("REPO_ROOT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("CACHE_ENABLED")
		os.Unsetenv("COPILOT_CLI_TIMEOUT")
	}()

	cfg := LoadFromEnv()

	if cfg.RepoRoot != "/custom/path" {
		t.Errorf("expected RepoRoot '/custom/path', got %q", cfg.RepoRoot)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got %q", cfg.LogLevel)
	}

	if cfg.CacheEnabled != false {
		t.Errorf("expected CacheEnabled false, got %v", cfg.CacheEnabled)
	}

	if cfg.CLITimeout != 60*time.Second {
		t.Errorf("expected CLITimeout 60s, got %v", cfg.CLITimeout)
	}
}

func TestLoadFromEnv_InvalidTimeout(t *testing.T) {
	os.Setenv("COPILOT_CLI_TIMEOUT", "invalid")
	defer os.Unsetenv("COPILOT_CLI_TIMEOUT")

	cfg := LoadFromEnv()

	// Should fall back to default
	if cfg.CLITimeout != 300*time.Second {
		t.Errorf("expected fallback to default 300s on invalid timeout, got %v", cfg.CLITimeout)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "uses environment value",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "uses default when not set",
			key:          "UNSET_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "parses true",
			key:          "BOOL_TRUE",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "parses false",
			key:          "BOOL_FALSE",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "uses default on invalid value",
			key:          "BOOL_INVALID",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
		},
		{
			name:         "uses default when not set",
			key:          "BOOL_UNSET",
			defaultValue: false,
			envValue:     "",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvBool(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		expected     time.Duration
	}{
		{
			name:         "parses valid duration",
			key:          "DUR_VALID",
			defaultValue: 10 * time.Second,
			envValue:     "30s",
			expected:     30 * time.Second,
		},
		{
			name:         "uses default on invalid duration",
			key:          "DUR_INVALID",
			defaultValue: 10 * time.Second,
			envValue:     "invalid",
			expected:     10 * time.Second,
		},
		{
			name:         "uses default when not set",
			key:          "DUR_UNSET",
			defaultValue: 15 * time.Second,
			envValue:     "",
			expected:     15 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvDuration(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
