# Abstractions - Coding Agent Instructions

## Project Overview

Abstractions is a boilerplate with a Go-based REST API and a Next.js frontend. The project is a **monorepo** with two independent applications that do not share code.

**Languages/Frameworks**: Go 1.26, TypeScript 5, Next.js 15 (App Router), React 19, Echo Framework v4
**Database**: PostgreSQL 15+ with the [pgvector](https://github.com/pgvector/pgvector) extension (used for semantic search/embeddings)
**Architecture**: Clean Architecture with clear separation between handlers, services, repositories

## Critical Rules - READ FIRST

1. **ALWAYS search** for existing implementations before writing new code. Consistency is paramount.
2. **NEVER run dev servers or builds** unless explicitly requested by the user.
3. **Ask before** introducing new dependencies or making significant architectural changes.
4. **Trust these instructions**. Only search if information here is incomplete or incorrect.
5. **NEVER create markdown files for examples or summary of your changes.** Document usage and flows using inline code comments.
6. **All API request and response fields must use camelCase** (e.g., `firstName`, `workspaceId`). This applies to all JSON fields in `api/server/requests/` and `api/server/responses/`.
7. **Icon sizing uses `text-*` classes, not `size-*` classes.** The `Icon` component renders Material Symbols, which are font-based. Icon size is controlled by `font-size` (via Tailwind `text-*` classes, e.g. `text-base`, `text-[2rem]`). The `size-*` classes (e.g. `size-8`) only set `width`/`height` and have no effect on the rendered icon size.
8. **Before creating a new component or writing custom code to solve a UI problem, always check `frontend/src/components/ui/` first.** A reusable component may already exist there (e.g. `Spinner`, `Icon`, `Button`, `Tooltip`). Prefer using existing components over writing new custom implementations.

## Build & Test

### Repo tooling

```bash
./setup.sh   # one-time: creates a local Python venv for pre-commit and registers the git hook
```

The pre-commit hook runs `scripts/run_api_gofmt.sh` on staged API Go files and `scripts/run_frontend_precommit.sh` (ESLint autofix + Prettier) on staged frontend files.

### API (Go)

**Location**: `api/` · **Go**: 1.26 (`api/go.mod`) · **DB**: PostgreSQL 15+ with pgvector

```bash
cd api
go mod tidy
go build -o build/main cmd/main.go
```

**Migrations** (`golang-migrate`, numbered files in `api/db/migrations/`):

```bash
go run cmd/main.go migrate 0          # run all pending
go run cmd/main.go migrate <N>        # up N steps (or down if negative)
go run cmd/main.go migrate force <version>   # fix a dirty migration state
```

**Tests** (require a Postgres test database):

```bash
cd api
chmod +x ./scripts/setup_test_db.sh && sh ./scripts/setup_test_db.sh   # drops/recreates {DATABASE_NAME}_test, enables pgvector, runs migrations
go test -v ./...
go test -v ./services/user   # single package
```

- Tests use `github.com/stretchr/testify/suite`; each suite embeds `testutils.TestSuite`, which wraps every test in a transaction rollback for isolation.
- Use `testutils.Seeder` (`s.seeder`) to create fixtures.

**Run the server**: `cd api && go run cmd/main.go` (requires `api/.env` — see below).

**Environment**: copy `api/.env.example` to `api/.env` and fill in real values. That file is the source of truth for required variables (core config, JWT/encryption secrets, Postgres, Redis, S3-compatible storage, email provider, OAuth, platform AI keys, optional Slack webhook).

### Frontend (Next.js)

**Location**: `frontend/` · **Node**: 22+ · **Package manager**: pnpm 9.15.4 (`packageManager` field pinned in `package.json`)

```bash
cd frontend
pnpm install
pnpm run dev     # ONLY when explicitly requested
pnpm run build
pnpm run type-check
pnpm run lint                # pnpm run lint:fix to autofix
pnpm run format              # pnpm run format:check to check only
```

**Tests**:

```bash
pnpm run test              # Jest unit tests
pnpm run test:coverage
pnpm run test:e2e          # Playwright — requires a running frontend and API server
```

**Environment**: `frontend/.env.local` needs `NEXT_PUBLIC_API_URL=http://localhost:1234`.

## Continuous Integration

`.github/workflows/tests.yml` runs on every push/PR to any branch:

- **test-api**: spins up `pgvector/pgvector:pg18-trixie` as a service container, writes `api/.env` with `DATABASE_URL`, runs `scripts/setup_test_db.sh`, then `go test -v ./...`. Go version comes from `api/go.mod` (`actions/setup-go` with `go-version-file`).
- **test-frontend**: pnpm 9.15.4 + Node 22, `pnpm install --frozen-lockfile`, then `pnpm run test:ci`.

`.github/workflows/build-deploy.yml` handles deployment — not relevant to local development.

To replicate the API job locally, run the same setup/migrate/test commands shown above under **API (Go) → Tests**.

## Project Architecture

### API

```
api/
├── cmd/main.go              # Entry point, handles migrate commands
├── app.go                   # Application bootstrap + startup health logging
├── config/                  # Configuration loading (config.go, environment.go, oauth.go)
├── db/
│   ├── connection.go        # GORM database connection
│   ├── middleware.go        # Transaction middleware for Echo
│   ├── migrate.go           # Migration runner
│   ├── migrations/          # SQL migration files (numbered, up/down pairs)
│   ├── models/              # GORM models
│   └── seeders/             # Idempotent DB seeders, run at startup
├── server/
│   ├── server.go
│   ├── routes/routes.go     # All route definitions
│   ├── handlers/            # One handler per domain, + base_handler.go
│   ├── middlewares/         # auth.go, services.go, etc.
│   ├── requests/            # Request validation structs (camelCase JSON)
│   └── responses/           # Response structs (camelCase JSON)
├── services/                # Business logic — one package per domain, service.go + *_test.go
├── repositories/            # Data access only, returns db/models
├── clients/                 # External services: email/ (Cloudflare, SES fallback), storage/ (S3/Minio)
├── utils/                   # One subdirectory per concern, named `{name}_utils/`
├── logger/                  # zerolog structured logging
└── tests/testutils/         # Shared test suite + seeder
```

**Key patterns**:

- Handlers validate via `requests/`, call services, return via `responses/`.
- Services own business logic; repositories are DB-access-only.
- Every request is wrapped in a transaction (`db.WrapInTransaction`); services are pulled from Echo context via `h.Services(c)`.
- Migrations are numbered sequentially with both up/down files — always create both, test both directions.

### Frontend (Next.js App Router)

```
frontend/src/
├── app/                     # Route groups: (auth)/, (dashboard)/, (content)/, plus god-mode/
├── components/
│   ├── ui/                  # shadcn/ui + custom extensions — check here before writing new UI
│   └── */                   # Feature components
├── store/                   # Redux Toolkit: modules/ (slices), services/ (RTK Query)
├── api/                     # Axios client layer, one file per domain
├── hooks/, lib/, types/
frontend/e2e/                # Playwright specs, fixtures/, utils/ (e.g. auth helpers)
```

**Key patterns**:

- Use the `useAPI` hook for component-level calls; axios methods from `src/api/` when called from Redux; thunks for flows that update global state.
- State: Redux Toolkit globally, React hooks locally. Styling: Tailwind CSS 4 + shadcn/ui. Forms: react-hook-form + yup/zod.
- Unit tests live in `src/components/ui/__tests__/`; no snapshot tests; use `@testing-library/user-event` over `fireEvent`.
- Import alias `@/` → `src/`. Never cross-import between `frontend/` and `api/`.

### Permissions

Role-based; defined in `services/authorization/`, enforced via `middlewares.RequirePermissions()`. Admin routes additionally require `AdminAuthMiddleware` + a god-mode flag. Permissions are workspace-scoped (e.g. `WorkspaceUpdate`, `WorkspaceManageUsers`).

## Common Development Patterns

**New API endpoint**: search `api/server/handlers/` for a similar one → add/update handler → request struct → response struct → service logic → route in `routes.go` → tests in `api/services/{domain}/*_test.go`.

**New DB table**: migration up/down files → GORM model in `db/models/` → repository → update `tests/testutils/seeder.go` → `go run cmd/main.go migrate 1`.

**New frontend page**: route under `src/app/` (reuse a route group for shared layout) → components → API calls in `src/api/{domain}.ts` → Redux slice if it needs global state.

## Before Committing

```bash
# API
cd api && go mod tidy && go build -o build/main cmd/main.go && sh ./scripts/setup_test_db.sh && go test -v ./...

# Frontend
cd frontend && pnpm install && pnpm run type-check && pnpm run lint && pnpm run test
```

For DB changes: test the migration both up and down, and update `tests/testutils/seeder.go` if fixtures are affected.

## Troubleshooting

- **Go build fails**: `go mod tidy`, confirm `go version` is 1.26+, confirm `api/.env` exists.
- **Tests fail with a DB connection error**: confirm Postgres is up (`pg_isready -h localhost -p 5432`), rerun `sh ./scripts/setup_test_db.sh`.
- **Migration stuck "dirty"**: `go run cmd/main.go migrate force <version>`, then check the `schema_migrations` table.
- **Frontend build fails**: `pnpm install`, confirm Node is 22+. `react-hooks/exhaustive-deps` warnings are expected and can be ignored.
- **Import errors**: frontend uses the `@/` alias for `src/`; never import across the `frontend/`/`api/` boundary.

## Additional Notes

- **Storage**: S3-compatible (AWS S3 or Minio).
- **Email**: Cloudflare Email Service by default; AWS SES as a fallback via `EMAIL_PROVIDER=ses`. Templates in `api/clients/email/`.
- **Sessions**: JWT-based auth with database-backed session storage.
- **OAuth**: Google and GitHub.
- **AI**: platform-default Anthropic/OpenAI/xAI keys are supported; workspaces can supply their own, and `DISABLE_DEFAULT_AI_USAGE=true` forces bring-your-own-key.
- **Deployments**: Docker support lives in `api/deploy/` — not used for local development.
