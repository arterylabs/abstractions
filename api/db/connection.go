package db

import (
	"abstractions/config"
	"abstractions/logger"
	"context"
	"fmt"
	"time"

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

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic(ctx, "failed to get underlying sql.DB", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db
}
