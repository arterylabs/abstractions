package config

import (
	"abstractions/logger"
	"context"
	"fmt"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ENV struct {
	ENV                   string `mapstructure:"ENV"`
	PORT                  string `mapstructure:"PORT"`
	SigningKey            string `mapstructure:"SIGNING_KEY"`
	TokenRefreshKey       string `mapstructure:"TOKEN_REFRESH_KEY"`
	APIURL                string `mapstructure:"API_URL"`
	ClientURL             string `mapstructure:"CLIENT_URL"`
	EncryptionKey32Bit    string `mapstructure:"ENCRYPTION_KEY_32_BIT"`
	GoogleClientID        string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret    string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	DatabaseHost          string `mapstructure:"DATABASE_HOST"`
	DatabaseUser          string `mapstructure:"DATABASE_USER"`
	DatabasePassword      string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName          string `mapstructure:"DATABASE_NAME"`
	DatabasePort          int    `mapstructure:"DATABASE_PORT"`
	DatabaseSSLMode       string `mapstructure:"DATABASE_SSL_MODE"`
	RedisHost             string `mapstructure:"REDIS_HOST"`
	RedisPort             int    `mapstructure:"REDIS_PORT"`
	AWSAccessKey          string `mapstructure:"AWS_ACCESS_KEY"`
	AWSSecretKey          string `mapstructure:"AWS_SECRET_KEY"`
	AWSRegion             string `mapstructure:"AWS_REGION"`
	AWSStorageEndpoint    string `mapstructure:"AWS_STORAGE_ENDPOINT"`
	AWSStorageBucket      string `mapstructure:"AWS_STORAGE_BUCKET"`
	DevelopmentEmail      string `mapstructure:"DEVELOPMENT_EMAIL"`
	GitHubClientID        string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret    string `mapstructure:"GITHUB_CLIENT_SECRET"`
	AnthropicAPIKey       string `mapstructure:"ANTHROPIC_API_KEY"`
	OpenAIAPIKey          string `mapstructure:"OPENAI_API_KEY"`
	DisableDefaultAIUsage bool   `mapstructure:"DISABLE_DEFAULT_AI_USAGE"`
}

var Environment *ENV

func LoadEnvironment() (*ENV, error) {
	// Load .env file

	// the .env file can exist in the current directory or a parent directory, read the env file in order
	// first in the current directory, then if not found in the parent directory

	ctx := context.Background()
	logger.Info(ctx, "loading environment...")

	// Try to load a .env file from common locations; not required when env vars are injected by the
	// container runtime (Docker, K8s, etc.).
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			if err := godotenv.Load("../../.env"); err != nil {
				logger.Warn(ctx, "no .env file found, relying on environment variables")
			}
		}
	}

	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// BindEnv is required for Unmarshal to read OS environment variables.
	// AutomaticEnv alone only works with explicit viper.Get() calls; without
	// BindEnv, Unmarshal sees an empty store when no .env file is loaded (e.g. Docker).
	t := reflect.TypeOf(ENV{})
	for i := 0; i < t.NumField(); i++ {
		if key, ok := t.Field(i).Tag.Lookup("mapstructure"); ok {
			viper.BindEnv(key)
		}
	}

	var env ENV

	if err := viper.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	Environment = &env
	return &env, nil
}

func (e *ENV) IsDevelopment() bool {
	return e.ENV == "development"
}
