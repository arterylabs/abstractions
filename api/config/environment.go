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
	ENV                      string `mapstructure:"ENV"`
	PORT                     string `mapstructure:"PORT"`
	SigningKey               string `mapstructure:"SIGNING_KEY"`
	TokenRefreshKey          string `mapstructure:"TOKEN_REFRESH_KEY"`
	APIURL                   string `mapstructure:"API_URL"`
	ClientURL                string `mapstructure:"CLIENT_URL"`
	EncryptionKey32Bit       string `mapstructure:"ENCRYPTION_KEY_32_BIT"`
	GoogleClientID           string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret       string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	DatabaseHost             string `mapstructure:"DATABASE_HOST"`
	DatabaseUser             string `mapstructure:"DATABASE_USER"`
	DatabasePassword         string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName             string `mapstructure:"DATABASE_NAME"`
	DatabasePort             int    `mapstructure:"DATABASE_PORT"`
	DatabaseSSLMode          string `mapstructure:"DATABASE_SSL_MODE"`
	RedisHost                string `mapstructure:"REDIS_HOST"`
	RedisPort                int    `mapstructure:"REDIS_PORT"`
	AWSAccessKey             string `mapstructure:"AWS_ACCESS_KEY"`
	AWSSecretKey             string `mapstructure:"AWS_SECRET_KEY"`
	AWSRegion                string `mapstructure:"AWS_REGION"`
	AWSStorageEndpoint       string `mapstructure:"AWS_STORAGE_ENDPOINT"`
	AWSStorageBucket         string `mapstructure:"AWS_STORAGE_BUCKET"`
	DevelopmentEmail         string `mapstructure:"DEVELOPMENT_EMAIL"`
	GitHubClientID           string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret       string `mapstructure:"GITHUB_CLIENT_SECRET"`
	AnthropicAPIKey          string `mapstructure:"ANTHROPIC_API_KEY"`
	OpenAIAPIKey             string `mapstructure:"OPENAI_API_KEY"`
	XaiAPIKey                string `mapstructure:"XAI_API_KEY"`
	DisableDefaultAIUsage    bool   `mapstructure:"DISABLE_DEFAULT_AI_USAGE"`
	SlackAppEventsWebhookURL string `mapstructure:"SLACK_APP_EVENTS_WEBHOOK_URL"`
	EmailProvider            string `mapstructure:"EMAIL_PROVIDER"`
	CloudflareAccountID      string `mapstructure:"CLOUDFLARE_ACCOUNT_ID"`
	CloudflareAPIToken       string `mapstructure:"CLOUDFLARE_API_TOKEN"`
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

// PublicConfigItem is a non-sensitive environment value safe to expose to god-mode admins.
type PublicConfigItem struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Category string `json:"category"`
}

// SecretConfigStatus reports whether a sensitive env var is set, without exposing its value.
type SecretConfigStatus struct {
	Key        string `json:"key"`
	Configured bool   `json:"configured"`
	Category   string `json:"category"`
}

// PublicAdminConfig is the sanitized view of loaded environment configuration.
type PublicAdminConfig struct {
	Config  []PublicConfigItem   `json:"config"`
	Secrets []SecretConfigStatus `json:"secrets"`
}

func configItem(category, key, value string) PublicConfigItem {
	return PublicConfigItem{Key: key, Value: value, Category: category}
}

func secretStatus(category, key, value string) SecretConfigStatus {
	return SecretConfigStatus{Key: key, Configured: value != "", Category: category}
}

// ToPublicAdminConfig returns non-sensitive env values and secret presence flags.
func (e *ENV) ToPublicAdminConfig() PublicAdminConfig {
	if e == nil {
		return PublicAdminConfig{
			Config:  []PublicConfigItem{},
			Secrets: []SecretConfigStatus{},
		}
	}

	boolStr := func(v bool) string {
		if v {
			return "true"
		}
		return "false"
	}

	return PublicAdminConfig{
		Config: []PublicConfigItem{
			configItem("General", "ENV", e.ENV),
			configItem("General", "PORT", e.PORT),
			configItem("General", "API_URL", e.APIURL),
			configItem("General", "CLIENT_URL", e.ClientURL),
			configItem("General", "DEVELOPMENT_EMAIL", e.DevelopmentEmail),
			configItem("General", "DISABLE_DEFAULT_AI_USAGE", boolStr(e.DisableDefaultAIUsage)),

			configItem("Database", "DATABASE_HOST", e.DatabaseHost),
			configItem("Database", "DATABASE_PORT", fmt.Sprintf("%d", e.DatabasePort)),
			configItem("Database", "DATABASE_NAME", e.DatabaseName),
			configItem("Database", "DATABASE_USER", e.DatabaseUser),
			configItem("Database", "DATABASE_SSL_MODE", e.DatabaseSSLMode),

			configItem("Redis", "REDIS_HOST", e.RedisHost),
			configItem("Redis", "REDIS_PORT", fmt.Sprintf("%d", e.RedisPort)),

			configItem("Storage", "AWS_REGION", e.AWSRegion),
			configItem("Storage", "AWS_STORAGE_ENDPOINT", e.AWSStorageEndpoint),
			configItem("Storage", "AWS_STORAGE_BUCKET", e.AWSStorageBucket),

			configItem("Email", "EMAIL_PROVIDER", e.EmailProvider),
			configItem("Email", "CLOUDFLARE_ACCOUNT_ID", e.CloudflareAccountID),

			configItem("OAuth", "GOOGLE_CLIENT_ID", e.GoogleClientID),
			configItem("OAuth", "GITHUB_CLIENT_ID", e.GitHubClientID),
		},
		Secrets: []SecretConfigStatus{
			secretStatus("Security", "SIGNING_KEY", e.SigningKey),
			secretStatus("Security", "TOKEN_REFRESH_KEY", e.TokenRefreshKey),
			secretStatus("Security", "ENCRYPTION_KEY_32_BIT", e.EncryptionKey32Bit),

			secretStatus("Database", "DATABASE_PASSWORD", e.DatabasePassword),

			secretStatus("Storage", "AWS_ACCESS_KEY", e.AWSAccessKey),
			secretStatus("Storage", "AWS_SECRET_KEY", e.AWSSecretKey),

			secretStatus("Email", "CLOUDFLARE_API_TOKEN", e.CloudflareAPIToken),

			secretStatus("OAuth", "GOOGLE_CLIENT_SECRET", e.GoogleClientSecret),
			secretStatus("OAuth", "GITHUB_CLIENT_SECRET", e.GitHubClientSecret),

			secretStatus("AI", "ANTHROPIC_API_KEY", e.AnthropicAPIKey),
			secretStatus("AI", "OPENAI_API_KEY", e.OpenAIAPIKey),
			secretStatus("AI", "XAI_API_KEY", e.XaiAPIKey),

			secretStatus("Integrations", "SLACK_APP_EVENTS_WEBHOOK_URL", e.SlackAppEventsWebhookURL),
		},
	}
}
