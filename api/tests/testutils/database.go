package testutils

import (
	"context"
	"fmt"

	"abstractions/logger"

	_ "github.com/lib/pq"
	gormdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testDBSuffix = "_test"
	testDBName   = "abstractions_test"
)

// SetupTestDatabase connects to the test database
func SetupTestDatabase(ctx context.Context) (*gorm.DB, error) {
	// Use local database connection parameters
	dbHost := "localhost"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbPort := 5432
	dbSSLMode := "disable"

	// Connect to test database
	testDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbHost, dbUser, dbPassword, testDBName, dbPort, dbSSLMode)

	// Create GORM DB instance
	gormDB, err := gorm.Open(gormdriver.Open(testDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM DB: %v", err)
	}

	logger.Info(ctx, "Connected to test database")
	return gormDB, nil
}

// CleanupTestDatabase doesn't need to do anything now since the database
// cleanup is handled by the setup script
func CleanupTestDatabase(ctx context.Context) error {
	return nil
}
