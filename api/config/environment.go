package config

import (
	"abstractions/logger"
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ENV struct {
	ENV                string `mapstructure:"ENV"`
	PORT               string `mapstructure:"PORT"`
	SigningKey         string `mapstructure:"SIGNING_KEY"`
	TokenRefreshKey    string `mapstructure:"TOKEN_REFRESH_KEY"`
	APIURL             string `mapstructure:"API_URL"`
	ClientURL          string `mapstructure:"CLIENT_URL"`
	EncryptionKey32Bit string `mapstructure:"ENCRYPTION_KEY_32_BIT"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	DatabaseHost       string `mapstructure:"DATABASE_HOST"`
	DatabaseUser       string `mapstructure:"DATABASE_USER"`
	DatabasePassword   string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName       string `mapstructure:"DATABASE_NAME"`
	DatabasePort       int    `mapstructure:"DATABASE_PORT"`
	DatabaseSSLMode    string `mapstructure:"DATABASE_SSL_MODE"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          int    `mapstructure:"REDIS_PORT"`
	AWSAccessKey       string `mapstructure:"AWS_ACCESS_KEY"`
	AWSSecretKey       string `mapstructure:"AWS_SECRET_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`
	AWSStorageEndpoint string `mapstructure:"AWS_STORAGE_ENDPOINT"`
	AWSStorageBucket   string `mapstructure:"AWS_STORAGE_BUCKET"`
	DevelopmentEmail   string `mapstructure:"DEVELOPMENT_EMAIL"`

	SlackClientID     string `mapstructure:"SLACK_CLIENT_ID"`
	SlackClientSecret string `mapstructure:"SLACK_CLIENT_SECRET"`

	GitHubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`

	AnthropicAPIKey string `mapstructure:"ANTHROPIC_API_KEY"`
	OpenAIAPIKey    string `mapstructure:"OPENAI_API_KEY"`
}

var Environment *ENV

func LoadEnvironment() (*ENV, error) {
	// Load .env file

	// the .env file can exist in the current directory or a parent directory, read the env file in order
	// first in the current directory, then if not found in the parent directory

	logger.Info(context.Background(), "loading environment...")
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			if err := godotenv.Load("../../.env"); err != nil {
				return nil, fmt.Errorf("error loading .env file: %w", err)
			}
		}
	}

	// if err := godotenv.Load("./.env"); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	viper.SetConfigFile("./.env")
	viper.SetEnvPrefix("")
	viper.ReadInConfig()

	var env ENV

	// Unmarshal environment variables into the struct
	if err := viper.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	Environment = &env
	return &env, nil
}

func (e *ENV) IsDevelopment() bool {
	return e.ENV == "development"
}
