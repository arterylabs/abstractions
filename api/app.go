package application

import (
	"context"
	"fmt"
	"log"

	"abstractions/config"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/routes"
)

func Start(cfg *config.Config) {
	if len(cfg.Environment.SigningKey) < 32 {
		log.Fatal("SIGNING_KEY must be at least 32 characters for HS256 security")
	}

	app := server.NewServer(cfg)

	routes.ConfigureRoutes(app)

	logStartupHealth(cfg)

	err := app.Start(cfg.Environment.PORT)
	if err != nil {
		log.Fatal("Port already used")
	}
}

// logStartupHealth prints non-sensitive configuration and connectivity status at startup.
func logStartupHealth(cfg *config.Config) {
	ctx := context.Background()
	env := cfg.Environment

	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logger.Info(ctx, "🚀 API starting up")
	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Core
	logger.Info(ctx, fmt.Sprintf("🌍  Environment : %s", env.ENV))
	logger.Info(ctx, fmt.Sprintf("🔌  Port        : %s", env.PORT))
	logger.Info(ctx, fmt.Sprintf("🌐  API URL     : %s", env.APIURL))
	logger.Info(ctx, fmt.Sprintf("🖥️   Client URL  : %s", env.ClientURL))

	// Database (no password)
	logger.Info(ctx, fmt.Sprintf("🗄️   Database    : %s@%s:%d/%s (sslmode=%s)",
		env.DatabaseUser, env.DatabaseHost, env.DatabasePort, env.DatabaseName, env.DatabaseSSLMode))

	// Redis
	logger.Info(ctx, fmt.Sprintf("🔴  Redis       : %s:%d", env.RedisHost, env.RedisPort))

	// AWS / Storage (no secret key)
	if env.AWSRegion != "" || env.AWSStorageEndpoint != "" {
		logger.Info(ctx, fmt.Sprintf("☁️   Storage     : region=%s endpoint=%s bucket=%s",
			env.AWSRegion, env.AWSStorageEndpoint, env.AWSStorageBucket))
	} else {
		logger.Warn(ctx, "⚠️   Storage     : AWS not configured")
	}

	// OAuth providers
	if env.GoogleClientID != "" {
		logger.Info(ctx, "✅  Google OAuth : configured")
	} else {
		logger.Warn(ctx, "⚠️   Google OAuth : not configured")
	}

	// Dev email (only relevant in non-production)
	if !config.IsProduction(cfg) && env.DevelopmentEmail != "" {
		logger.Info(ctx, fmt.Sprintf("📧  Dev email    : %s", env.DevelopmentEmail))
	}

	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
