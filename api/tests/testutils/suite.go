package testutils

import (
	"abstractions/config"
	"abstractions/server"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// TestSuite is a base test suite that provides common functionality
type TestSuite struct {
	suite.Suite
	DB         *gorm.DB
	Server     *server.Server
	Ctx        context.Context
	Config     *config.Config
	dbLockHeld bool
}

var testDBMutex sync.Mutex

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
	testDBMutex.Lock()
	s.dbLockHeld = true

	// Get all tables in the public schema
	var tables []string
	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public'
		  AND tablename <> 'schema_migrations'
		ORDER BY tablename
	`
	if err := s.DB.Raw(query).Scan(&tables).Error; err != nil {
		s.T().Errorf("Failed to get table names: %v", err)
		return
	}

	if len(tables) == 0 {
		return
	}

	quotedTables := make([]string, 0, len(tables))
	for _, table := range tables {
		quotedTables = append(quotedTables, fmt.Sprintf("%q", table))
	}

	// Truncate in a single statement so each test starts from a consistent state
	// without cascading lock contention between per-table TRUNCATE calls.
	if err := s.DB.Exec("TRUNCATE TABLE " + strings.Join(quotedTables, ", ") + " CASCADE").Error; err != nil {
		s.T().Errorf("Failed to truncate test tables: %v", err)
	}
}

// TearDownTest runs after each test
func (s *TestSuite) TearDownTest() {
	if s.dbLockHeld {
		s.dbLockHeld = false
		testDBMutex.Unlock()
	}
}

// RunTestSuite is a helper function to run the test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
