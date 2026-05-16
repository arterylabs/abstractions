# Abstractions - Copilot Coding Instructions

## Project Overview

Abstractions is a boilerplate with a Go-based REST API and a Next.js frontend. The project is a **monorepo** with two independent applications that do not share code.

**Languages/Frameworks**: Go 1.23+, TypeScript 5, Next.js 15 (App Router), React 19, Echo Framework v4  
**Architecture**: Clean Architecture with clear separation between handlers, services, repositories

## Critical Rules - READ FIRST

1. **ALWAYS search** for existing implementations before writing new code. Consistency is paramount.
2. **NEVER run dev servers or builds** unless explicitly requested by the user.
3. **Ask before** introducing new dependencies or making significant architectural changes.
4. **Trust these instructions**. Only search if information here is incomplete or incorrect.
5. **NEVER create markdown files for examples or summary of your changes** document usage and flows using inline code comments.
6. **All API request and response fields must use camelCase** (e.g., `firstName`, `workspaceId`). This applies to all JSON fields in `api/server/requests/` and `api/server/responses/`.
7. **Icon sizing uses `text-*` classes, not `size-*` classes.** The `Icon` component renders Material Symbols, which are font-based. Icon size is controlled by `font-size` (via Tailwind `text-*` classes, e.g. `text-base`, `text-[2rem]`). The `size-*` classes (e.g. `size-8`) only set `width`/`height` and have no effect on the rendered icon size.
8. **Before creating a new component or writing custom code to solve a UI problem, always check `frontend/src/components/ui/` first.** A reusable component may already exist there (e.g. `Spinner`, `Icon`, `Button`, `Tooltip`). Prefer using existing components over writing new custom implementations.

## Build & Test Instructions

### Repo Tooling

```bash
# Install pre-commit and register the git hook
./setup.sh

# Run gofmt across all Go files in api/
sh api/scripts/gofmt.sh
```

- `setup.sh` creates a repo-local Python virtualenv for `pre-commit` and registers the repository's pre-commit hook.
- The pre-commit hook runs `gofmt` on staged API Go files plus ESLint autofix and Prettier formatting on staged frontend files.

### API (Go Backend)

**Location**: `/api` directory  
**Go Version**: 1.23.4 (specified in go.mod toolchain)  
**Database**: PostgreSQL 15

#### Prerequisites

- PostgreSQL 15 running on localhost:5432
- Go 1.23+ installed
- Required environment variables in `.env` file (see Environment Setup)

#### Build Commands

```bash
cd api
go mod tidy                          # Always run first to sync dependencies
go build -o build/main cmd/main.go   # Build binary to build/main
```

#### Running Migrations

```bash
cd api
# Run all pending migrations
go run cmd/main.go migrate 0

# Run specific number of steps
go run cmd/main.go migrate <N>       # Up N steps (positive) or down N steps (negative)

# Force migration version (fix dirty state)
go run cmd/main.go migrate force <version>

# With custom database (useful for tests)
go run cmd/main.go migrate \
  --host localhost \
  --port 5432 \
  --user postgres \
  --password postgres \
  --dbname abstractions_test \
  --sslmode disable \
  0
```

#### Testing

**CRITICAL**: Tests require a PostgreSQL test database setup.

```bash
cd api

# Setup test database (must run first)
chmod +x ./scripts/setup_test_db.sh
sh ./scripts/setup_test_db.sh

# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./services/user
```

**Test Database Setup**:

- Database: `abstractions_test` (automatically created by setup script)
- Connection: `postgresql://postgres:postgres@localhost:5432/abstractions_test?sslmode=disable`
- The setup script drops/recreates the database and runs migrations
- Tests use `testutils.TestSuite` which provides transaction isolation

**Known Test Patterns**:

- Tests use `github.com/stretchr/testify/suite` for test organization
- Each test suite embeds `testutils.TestSuite` for DB setup
- Use `s.seeder` (from `testutils.Seeder`) to create test data

#### Running the API Server

```bash
cd api
go run cmd/main.go
# Requires .env file with all configuration (see Environment Setup)
```

#### Environment Setup (.env file required in /api directory)

The API requires a `.env` file in the `api/` directory. Key variables:

```
ENV=development
PORT=1234
SIGNING_KEY=<jwt_secret>
TOKEN_REFRESH_KEY=<refresh_secret>
API_URL=http://localhost:1234
CLIENT_URL=http://localhost:3000
ENCRYPTION_KEY_32_BIT=<32_char_key>
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=abstractions
DATABASE_PORT=5432
DATABASE_SSL_MODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
AWS_ACCESS_KEY=<aws_key>
AWS_SECRET_KEY=<aws_secret>
AWS_REGION=<region>
AWS_STORAGE_ENDPOINT=<s3_endpoint>
AWS_STORAGE_BUCKET=<bucket>
GOOGLE_CLIENT_ID=<oauth_client_id>
GOOGLE_CLIENT_SECRET=<oauth_secret>
DEVELOPMENT_EMAIL=<email>
```

### Frontend (Next.js)

**Location**: `/frontend` directory  
**Node Version**: v22.12.0  
**Package Manager**: pnpm 9.15.4  
**Framework**: Next.js 15

#### Install Dependencies

```bash
cd frontend
pnpm install    # Always run after pulling changes or adding dependencies
```

#### Build Commands

```bash
cd frontend

# Type checking (runs during build, can run separately)
pnpm run type-check

# Linting (runs during build, warnings are acceptable)
pnpm run lint

# Autofix lint issues
pnpm run lint:fix

# Format files with Prettier
pnpm run format

# Check formatting without writing changes
pnpm run format:check

# Production build
pnpm run build

# Development mode (ONLY when explicitly requested)
pnpm run dev --turbopack

# Production preview
pnpm run start
```

#### Testing

```bash
cd frontend

# Run Jest unit tests
pnpm run test

# Run with coverage
pnpm run test:coverage

# Run e2e tests (requires a running frontend server)
pnpm run test:e2e
```

#### Environment Variables

Frontend requires `.env.local` file in `frontend/` directory:

```
NEXT_PUBLIC_API_URL=http://localhost:1234
```

## Continuous Integration

### GitHub Actions Workflow

**File**: `.github/workflows/tests.yml`  
**Trigger**: All pushes and pull requests to any branch  
**Jobs**: RunAPIUnitTests + RunFrontendUnitTests

**Steps**:

1. Checkout code
2. Setup Go 1.23 with caching
3. Install dependencies: `cd api && go mod tidy && go get -d ./...`
4. Create `.env` file with test database URL
5. Setup PostgreSQL service (runs in Docker)
6. Run test database setup script: `sh ./scripts/setup_test_db.sh`
7. Run tests: `cd api && go test -v ./...`
8. Run frontend unit tests: `cd frontend && pnpm run test:ci`

**To replicate CI locally**:

```bash
cd api
go mod tidy
go get -d ./...
echo "DATABASE_URL=postgresql://postgres:postgres@localhost:5432/abstractions_test?sslmode=disable" > .env
sh ./scripts/setup_test_db.sh
go test -v ./...
```

## Project Architecture

### API Architecture (Clean Architecture)

```
api/
├── cmd/main.go              # Entry point, handles migrate commands
├── app.go                   # Application bootstrap + startup health logging
├── config/                  # Configuration loading
│   ├── config.go           # Main config struct
│   ├── environment.go      # ENV struct and loading from .env
│   └── oauth.go            # OAuth provider configs
├── db/
│   ├── connection.go       # GORM database connection
│   ├── middleware.go       # Transaction middleware for Echo
│   ├── migrate.go          # Migration runner
│   ├── migrations/         # SQL migration files (numbered)
│   ├── models/             # GORM models (database entities)
│   └── seeders/            # Database seeders (run at startup)
├── server/
│   ├── server.go           # Server struct and initialization
│   ├── routes/routes.go    # All route definitions
│   ├── handlers/           # HTTP request handlers
│   │   ├── *_handler.go   # One handler per domain
│   │   └── base_handler.go # Base handler with common methods
│   ├── middlewares/        # Echo middlewares
│   │   ├── auth.go        # Authentication middleware
│   │   └── services.go    # Services injection middleware
│   ├── requests/           # Request validation structs
│   └── responses/          # Response formatting structs
├── services/               # Business logic layer
│   ├── services.go         # Services aggregate struct
│   └── */                  # One package per domain
│       ├── service.go      # Service implementation
│       └── *_test.go       # Service tests
├── repositories/           # Data access layer (database queries)
│   └── *_repository.go
├── clients/                # External service clients
│   ├── email/             # Email service (SES)
│   └── storage/           # File storage (S3/Minio)
├── utils/                  # Utility packages
│   └── */                 # Organized by concern (all in subdirectories with _utils suffix)
├── logger/                 # Structured logging (zerolog)
└── tests/
    └── testutils/          # Test utilities and base suite
```

**Key Patterns**:

- **Handlers** receive HTTP requests, validate using `requests/` structs, call services, return `responses/`
- **Services** contain business logic, orchestrate repositories and external clients
- **Repositories** perform database operations only, return `db/models`
- **Middleware** wraps each request in a transaction (`db.WrapInTransaction`) and injects services
- **Context propagation**: Services are accessed via `h.Services(c)` which pulls from Echo context
- **Migrations**: Numbered sequentially (000001*, 000002*, etc.) with up/down SQL files
- **Seeders**: Safe to run multiple times (idempotent), executed on server startup
- **Utils naming**: All utility packages live in subdirectories named `{name}_utils/` (e.g. `jwt_utils/`, `ip_utils/`)

### Frontend Architecture (Next.js App Router)

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (auth)/            # Auth route group
│   │   │   ├── login/
│   │   │   ├── signup/
│   │   │   └── verify/
│   │   ├── (dashboard)/       # Dashboard route group
│   │   ├── (content)/         # Content pages
│   │   ├── god-mode/          # Admin panel
│   │   ├── layout.tsx         # Root layout
│   │   └── globals.css        # Global styles (Tailwind)
│   ├── components/
│   │   ├── ui/                # shadcn/ui components + custom extensions
│   │   │   └── __tests__/    # Unit tests for every UI component
│   │   └── */                 # Feature components
│   ├── store/                 # Redux Toolkit
│   │   ├── index.ts          # Store configuration
│   │   ├── hooks.ts          # Typed hooks (useAppDispatch, useAppSelector)
│   │   ├── modules/          # Redux slices
│   │   └── services/         # RTK Query services
│   ├── api/                   # API client layer
│   │   ├── index.ts          # Axios instance configuration
│   │   └── *.ts              # API endpoint functions
│   ├── hooks/                 # Custom React hooks
│   ├── lib/                   # Utility functions
│   └── types/                 # TypeScript types
├── e2e/                       # Playwright e2e tests
│   ├── fixtures/             # Shared test fixtures
│   ├── utils/                # Test helpers (auth, etc.)
│   └── */                    # Feature test specs
├── public/                    # Static assets
├── components.json            # shadcn/ui configuration
├── jest.config.js             # Jest unit test configuration
├── playwright.config.ts       # Playwright e2e configuration
├── tsconfig.json             # TypeScript strict mode
├── eslint.config.mjs         # ESLint (Next.js + TypeScript)
├── next.config.ts            # Next.js config with MDX support
└── package.json              # Dependencies and scripts
```

**Key Patterns**:

- **API Calls**:
  - Use `useAPI` hook for component-level API calls
  - Use axios methods from `src/api/` when called from Redux store
  - Redux thunk actions for operations that update central state
- **State Management**: Redux Toolkit for global state, React hooks for local state
- **Styling**: Tailwind CSS 4 with shadcn/ui components
- **Type Safety**: TypeScript strict mode enabled
- **Forms**: react-hook-form with yup/zod validation
- **Unit Tests**: Jest + React Testing Library; tests live in `src/components/ui/__tests__/`
- **E2E Tests**: Playwright; tests live in `e2e/`

### Permission System

The API uses a role-based permission system:

- Permissions defined in `services/authorization/`
- Checked via `middlewares.RequirePermissions()` middleware on routes
- Admin routes require `AdminAuthMiddleware` + god mode flag
- Workspace-level permissions (e.g., `WorkspaceUpdate`, `WorkspaceManageUsers`)

### Database Migrations

- Located in `api/db/migrations/`
- Named: `{number}_{description}.up.sql` and `{number}_{description}.down.sql`
- Managed by golang-migrate library
- Always create both up and down migrations
- Test migrations in both directions before committing

## Common Development Patterns

### Adding a New API Endpoint

1. Search for similar endpoints in `api/server/handlers/`
2. Create/update handler method following existing patterns
3. Add request validation struct in `api/server/requests/`
4. Add response struct in `api/server/responses/`
5. Implement service logic in `api/services/{domain}/`
6. Add route in `api/server/routes/routes.go`
7. Add tests in `api/services/{domain}/*_test.go`

### Adding a New Database Table

1. Create migration files: `api/db/migrations/{number}_{name}.up.sql` and `.down.sql`
2. Create GORM model in `api/db/models/{name}.go`
3. Create repository in `api/repositories/{name}_repository.go`
4. Update test seeder in `api/tests/testutils/seeder.go`
5. Run migration: `go run cmd/main.go migrate 1`

### Adding a Frontend Page

1. Create route in `src/app/` following App Router conventions
2. Use route groups for layout sharing: `(auth)`, `(dashboard)`
3. Create components in `src/components/`
4. Add API calls to `src/api/{domain}.ts`
5. Add Redux slice if managing global state in `src/store/modules/`

## Validation Steps Before Committing

1. **API Changes**:

   ```bash
   cd api
   go mod tidy
   go build -o build/main cmd/main.go
   sh ./scripts/setup_test_db.sh
   go test -v ./...
   ```

2. **Frontend Changes**:

   ```bash
   cd frontend
   pnpm install
   pnpm run type-check
   pnpm run lint
   pnpm run test
   ```

3. **Database Changes**:
   - Test migration up and down
   - Update test data seeders
   - Run full test suite

## Troubleshooting Common Issues

### Go Build Fails

- Run `go mod tidy` first
- Check Go version: `go version` (need 1.23+)
- Ensure `api/.env` file exists

### Tests Fail with Database Connection Error

- Ensure PostgreSQL is running: `pg_isready -h localhost -p 5432`
- Run test database setup: `sh ./scripts/setup_test_db.sh`
- Check test database exists: `psql -U postgres -l | grep abstractions_test`

### Migration "Dirty" State

- Force migration version: `go run cmd/main.go migrate force <version>`
- Manually check schema_migrations table in database

### Frontend Build Fails

- Run `pnpm install` first
- Check Node version: `node --version` (need v22+)
- Ignore react-hooks/exhaustive-deps warnings (expected)

### Import Path Issues

- Frontend uses `@/` alias for `src/` directory
- Never import from `frontend/` in frontend code
- Never import from `api/` in frontend code

## Additional Notes

- **Storage**: File uploads use S3-compatible storage (AWS S3 or Minio)
- **Email**: SES for transactional emails, templates in `api/clients/email/`
- **Logging**: Structured logging with zerolog throughout API
- **Sessions**: JWT-based authentication with database session storage
- **OAuth**: Google and Github OAuth integrations supported
- **Deployments**: Docker support in `api/deploy/` (not for local development)
