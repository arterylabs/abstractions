package testutils

import (
	"abstractions/config"
	"abstractions/server"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// testDBAdvisoryLockKey serializes access to the shared abstractions_test database
// across package processes. go test runs packages in parallel by default; an
// in-process mutex alone cannot prevent concurrent TRUNCATE / fixture races.
// Value is an arbitrary stable int64 ("Abst" in ASCII).
const testDBAdvisoryLockKey int64 = 0x41627374

// TestSuite is a base test suite that provides common functionality
type TestSuite struct {
	suite.Suite
	DB         *gorm.DB
	Server     *server.Server
	Ctx        context.Context
	Config     *config.Config
	dbLockHeld bool
	// lockConn holds a dedicated connection that owns the session-level
	// advisory lock for the duration of a single test.
	lockConn *sql.Conn
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

// acquireTestDBLock takes an in-process mutex and a PostgreSQL session-level
// advisory lock so only one test (across all packages/processes) mutates the
// shared test database at a time.
func (s *TestSuite) acquireTestDBLock() {
	testDBMutex.Lock()
	s.dbLockHeld = true

	sqlDB, err := s.DB.DB()
	if err != nil {
		s.T().Fatalf("Failed to get underlying sql.DB for test lock: %v", err)
	}

	conn, err := sqlDB.Conn(s.Ctx)
	if err != nil {
		s.T().Fatalf("Failed to open dedicated connection for test lock: %v", err)
	}

	if _, err := conn.ExecContext(s.Ctx, "SELECT pg_advisory_lock($1)", testDBAdvisoryLockKey); err != nil {
		_ = conn.Close()
		s.T().Fatalf("Failed to acquire test database advisory lock: %v", err)
	}
	s.lockConn = conn
}

// releaseTestDBLock drops the advisory lock and the in-process mutex.
func (s *TestSuite) releaseTestDBLock() {
	if s.lockConn != nil {
		// Best-effort unlock; closing the connection also releases session locks.
		_, _ = s.lockConn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", testDBAdvisoryLockKey)
		_ = s.lockConn.Close()
		s.lockConn = nil
	}
	if s.dbLockHeld {
		s.dbLockHeld = false
		testDBMutex.Unlock()
	}
}

// SetupTest runs before each test
func (s *TestSuite) SetupTest() {
	s.acquireTestDBLock()

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
	s.releaseTestDBLock()
}

// RunTestSuite is a helper function to run the test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
