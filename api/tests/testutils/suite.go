package testutils

import (
	"abstractions/config"
	"abstractions/server"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// TestSuite is a base test suite that provides common functionality
type TestSuite struct {
	suite.Suite
	DB     *gorm.DB
	Server *server.Server
	Ctx    context.Context
	Config *config.Config
}

// SetupSuite runs once before all tests in the suite
func (s *TestSuite) SetupSuite() {
	s.Ctx = context.Background()

	// Connect to test database and run migrations
	var err error
	s.DB, err = SetupTestDatabase(s.Ctx)
	if err != nil {
		s.T().Fatalf("Failed to setup test database: %v", err)
	} else {
		s.T().Log("Connected to test database")
	}

	s.Config = TestConfig()

	// Set the global config.Environment for services that depend on it
	config.Environment = s.Config.Environment

	// Use existing TestConfig from config.go
	s.Server = server.NewServer(s.Config)
}

// TearDownSuite runs once after all tests in the suite
func (s *TestSuite) TearDownSuite() {
	if err := CleanupTestDatabase(s.Ctx); err != nil {
		s.T().Errorf("Failed to cleanup test database: %v", err)
	}
}

// SetupTest runs before each test
func (s *TestSuite) SetupTest() {
	// Get all tables in the public schema
	var tables []string
	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public'
	`
	if err := s.DB.Raw(query).Scan(&tables).Error; err != nil {
		s.T().Errorf("Failed to get table names: %v", err)
		return
	}

	// Start each test with a clean slate by truncating all tables
	for _, table := range tables {
		if err := s.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			s.T().Errorf("Failed to truncate table %s: %v", table, err)
		}
	}
}

// TearDownTest runs after each test
func (s *TestSuite) TearDownTest() {
	// Clean up any test data if needed
}

// RunTestSuite is a helper function to run the test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
