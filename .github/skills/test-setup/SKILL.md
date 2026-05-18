---
name: test-setup
description: Documents the Abstractions API testing infrastructure, patterns, environment setup, and how to run tests. Use this when asked about how tests work, the test setup, test patterns, or how to run the test suite.
---

You are an expert on the Abstractions API test infrastructure. Use the context below to answer questions accurately and help developers set up, run, and write tests correctly.

---

## Overview

The Abstractions test suite covers the **Go API** (`api/` directory) with integration-style tests against a real PostgreSQL test database, plus **frontend** tests: Jest + React Testing Library unit tests (`frontend/src/`) and Playwright end-to-end tests (`frontend/e2e/`).

---

## Test Database Setup

The test database must be provisioned before running tests.

```bash
cd api
chmod +x ./scripts/setup_test_db.sh
sh ./scripts/setup_test_db.sh
```

This script:
1. Drops and recreates the `harmony_test` database on `localhost:5432`
2. Runs all migrations against it using `go run cmd/main.go migrate`

Connection string used by tests:
```
postgresql://postgres:postgres@localhost:5432/harmony_test?sslmode=disable
```

Tests pick this up via a `.env` file in `api/` containing:
```
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/harmony_test?sslmode=disable
```

---

## Running Tests

```bash
# All tests
cd api && go test -v ./...

# Single package
cd api && go test -v ./services/event

# With race detector
cd api && go test -race ./...
```

---

## Test Infrastructure

### `testutils.TestSuite` (`api/tests/testutils/suite.go`)

The base struct that all test suites embed. Provides:

| Field | Type | Purpose |
|-------|------|---------|
| `s.DB` | `*gorm.DB` | Per-test transaction — use for all DB access |
| `s.Ctx` | `context.Context` | Background context |
| `s.Config` | `*config.Config` | Test configuration |
| `s.Server` | `*server.Server` | Initialized server (for middleware tests) |

**Isolation strategy**:
- `SetupSuite` — runs once; deletes all committed rows from the test DB before the suite begins
- `SetupTest` — opens a fresh `BEGIN` transaction for every individual test
- `TearDownTest` — `ROLLBACK`s the transaction; writes are never committed, so concurrent packages cannot interfere

### `testutils.Seeder` (`api/tests/testutils/seeder.go`)

Factory for creating test fixtures inside the per-test transaction.

```go
seeder := testutils.NewSeeder(s.T(), s.DB)

user  := seeder.CreateUser()
event := seeder.CreateEvent(user)
guest := seeder.CreateGuest(event, user)
log   := seeder.CreateActivityLog(event, user, models.ActivityTypeGuestCreated, nil)
```

Available factory methods:
- `CreateUser(opts ...func(*models.User)) *models.User`
- `CreateEvent(user, opts ...) *models.Event`
- `CreateUserEventRelationship(user, event, role)`
- `CreateAuthenticator(user, type, token)`
- `CreateGuest(event, user, opts ...)`
- `CreateActivityLog(event, user, type, metadata)`

Option helpers: `WithEmail(email)`, `WithUserName(name)`, `WithEventName(name)`

---

## Writing a New Test Suite

Follow this exact pattern, used consistently across all service tests:

```go
package myfeature_service_test

import (
    "testing"
    myfeature_service "abstractions/services/myfeature"
    "abstractions/tests/testutils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type MyFeatureTestSuite struct {
    testutils.TestSuite
    seeder    *testutils.Seeder
    myService *myfeature_service.Service
}

func TestMyFeatureService(t *testing.T) {
    suite.Run(t, new(MyFeatureTestSuite))
}

func (s *MyFeatureTestSuite) SetupTest() {
    s.TestSuite.SetupTest() // MUST call this first
    s.seeder = testutils.NewSeeder(s.T(), s.DB)
    s.myService = myfeature_service.NewService(s.DB)
}

func (s *MyFeatureTestSuite) TestSomeBehavior() {
    user := s.seeder.CreateUser()
    // ... use s.Ctx and s.DB in assertions
    result, err := s.myService.DoSomething(s.Ctx, user.ID)
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), result)
}
```

**Critical rules**:
- Always call `s.TestSuite.SetupTest()` first in `SetupTest()`
- Always use `s.DB` (not a new connection) for all service constructors
- Always use `s.Ctx` as the context argument
- Use `assert.*` (not `require.*`) unless a nil pointer would panic subsequent lines

---

## Test Locations

| Layer | Location | Pattern |
|-------|----------|---------|
| Services | `api/services/{domain}/*_test.go` | Integration against test DB |
| Repositories | `api/repositories/enum_validation_test.go` | Enum validation only |
| Handlers/Middlewares | `api/server/middlewares/*_test.go` | Unit tests |
| Utilities | `api/utils/{pkg}/*_test.go` | Unit tests |
| Email client | `api/clients/email/email_test.go` | Unit tests |
| Logger | `api/logger/logger_test.go` | Unit tests |

---

## Known Skipped Tests

Tests that depend on external services are either skipped or run in a no-op mode:
- **S3/Minio storage** — file upload tests skip actual storage; `TestUploadFileSuccess` uses a mock
- **AWS SES email** — `SendVerificationEmail` is called but does not assert delivery
- **Twilio / WhatsApp SMS** — no test file exists; HTTP calls would need mocking

---

## CI Configuration

GitHub Actions workflow: `.github/workflows/tests.yml`

Steps mirror local setup:
1. Setup Go 1.23
2. `cd api && go mod tidy && go get -d ./...`
3. Start PostgreSQL 15 service container
4. Write `.env` with `DATABASE_URL`
5. `sh ./scripts/setup_test_db.sh`
6. `cd api && go test -v ./...`
