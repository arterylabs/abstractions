package testutils

import (
	"abstractions/config"
	"os"
	"strconv"
)

// TestConfig returns a configuration for testing
func TestConfig() *config.Config {
	cfg := &config.Config{
		Environment: &config.ENV{},
	}

	// Set test database configuration
	cfg.Environment.DatabaseHost = os.Getenv("TEST_DB_HOST")
	if cfg.Environment.DatabaseHost == "" {
		cfg.Environment.DatabaseHost = "localhost"
	}
	cfg.Environment.DatabasePort = 5432 // Default port as int
	if port := os.Getenv("TEST_DB_PORT"); port != "" {
		// Try to parse port as int, fallback to 5432 if invalid
		if portNum, err := strconv.Atoi(port); err == nil {
			cfg.Environment.DatabasePort = portNum
		}
	}
	cfg.Environment.DatabaseUser = os.Getenv("TEST_DB_USER")
	if cfg.Environment.DatabaseUser == "" {
		cfg.Environment.DatabaseUser = "postgres"
	}
	cfg.Environment.DatabasePassword = os.Getenv("TEST_DB_PASSWORD")
	if cfg.Environment.DatabasePassword == "" {
		cfg.Environment.DatabasePassword = "postgres"
	}
	cfg.Environment.DatabaseName = os.Getenv("TEST_DB_NAME")
	if cfg.Environment.DatabaseName == "" {
		cfg.Environment.DatabaseName = "abstractions_test"
	}
	cfg.Environment.DatabaseSSLMode = "disable" // Commonly used for tests

	// Set test encryption key (32 characters for AES-256)
	cfg.Environment.EncryptionKey32Bit = os.Getenv("TEST_ENCRYPTION_KEY_32_BIT")
	if cfg.Environment.EncryptionKey32Bit == "" {
		cfg.Environment.EncryptionKey32Bit = "12345678901234567890123456789012" // 32-character test key
	}

	// Set other required environment variables for tests
	cfg.Environment.ENV = "test"
	cfg.Environment.PORT = "8080"
	cfg.Environment.SigningKey = "test_signing_key_for_jwt_tokens_123456"
	cfg.Environment.TokenRefreshKey = "test_refresh_key_for_jwt_tokens_123456"
	cfg.Environment.APIURL = "http://localhost:8080"
	cfg.Environment.ClientURL = "http://localhost:3000"

	return cfg
}
