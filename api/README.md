# API

Go REST API, background workers, and database migrations for Abstractions.

See the [root README](../README.md) for product overview, full local setup, environment variables, and self-hosting notes.

**License:** [GNU Affero General Public License v3.0](../LICENSE) (AGPL-3.0) — same as the monorepo root.

## Quick commands

```bash
# From this directory
cp .env.example .env   # then edit values

go mod tidy
go run cmd/main.go migrate 0
go run cmd/main.go

# Workers (separate process)
go run workers/workers.go

# Tests (requires Postgres test DB)
sh ./scripts/setup_test_db.sh
go test -v ./...
```

## Layout

| Path | Role |
|------|------|
| `cmd/` | API entrypoint |
| `server/` | HTTP handlers, routes, middleware |
| `services/` | Business logic |
| `repositories/` | Database access |
| `workers/` | Async jobs (indexing, newsletters) |
| `db/migrations/` | SQL migrations |
| `clients/` | Email, storage, AI, GitHub, Redis |
