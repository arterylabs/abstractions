package db

import (
	"abstractions/config"
	"abstractions/logger"
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, cfg *config.Config) *gorm.DB {
	// Define the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Environment.DatabaseHost, cfg.Environment.DatabaseUser, cfg.Environment.DatabasePassword, cfg.Environment.DatabaseName, cfg.Environment.DatabasePort, cfg.Environment.DatabaseSSLMode)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Handle any errors that occur during the connection
		logger.Panic(ctx, "failed to connect to database ", err)
	}

	return db
}
