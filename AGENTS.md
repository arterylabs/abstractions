# Abstractions — Agent Instructions

Monorepo: Go REST API (`api/`) + Next.js frontend (`frontend/`). They do not share code. Stack: Go 1.26, Echo v4, PostgreSQL + pgvector, Next.js 15 App Router, React 19, pnpm, TypeScript 5.

## Rules

1. **Search first** — match existing patterns before writing new code.
2. **No unsolicited builds/dev servers** — only when the user asks.
3. **Ask before** new dependencies or architectural changes.
4. **No markdown docs** for examples or change summaries — use inline comments.
5. **API JSON is camelCase** — all fields in `api/server/requests/` and `api/server/responses/`.
6. **Icons**: size with `text-*` (font-size), not `size-*`. Prefer components in `frontend/src/components/ui/` over new UI.
7. **Never cross-import** between `frontend/` and `api/`. Frontend alias: `@/` → `src/`.
8. **Tests protect product confidence, not coverage %.** Before adding a test, answer: what user-visible or security regression does this prevent? If none, do not write it.

## Testing policy

- **Fewer tests are better.** Prefer one meaningful integration/service test over many implementation unit tests.
- **Do not unit-test** shadcn/Radix wrappers under `frontend/src/components/ui/` unless they contain **custom business logic** (validation, transforms, destructive confirmation gates, polling contracts). No “renders without crashing,” CSS class, label, or icon-existence tests.
- **API unit/integration tests** should cover business rules, auth/permissions, persistence side effects, scheduling/retries, encryption, and security boundaries — not constructors, constant maps, trivial getters, exact error message strings, or library wrappers (e.g. thin Redis clients).
- **E2E tests are user journeys** (signup, login, invite accept, connect repo, destructive gates). Do not assert marketing copy, page presence, empty-state wording, or tab labels. Prefer outcomes (URL, persisted data, can/cannot act).
- **Avoid retesting the same path** at repository + service + middleware layers; keep the highest-signal layer.
- **Mock only when isolation is required.** Prefer real DB transaction suites (`testutils.TestSuite`) for API behavior.

## Commands

```bash
# One-time repo setup (pre-commit hooks)
./setup.sh

# API
cd api
go mod tidy
go build -o build/main cmd/main.go
go run cmd/main.go                          # needs api/.env (from .env.example)
go run cmd/main.go migrate 0                # all pending; N = steps; force <v> if dirty
sh ./scripts/setup_test_db.sh && go test -v ./...

# Frontend (Node 22+, pnpm 9.15.4)
cd frontend
pnpm install
pnpm run type-check && pnpm run lint && pnpm run test
# pnpm run dev / build / test:e2e — only when requested
```

Env: `api/.env` from `api/.env.example`; `frontend/.env.local` needs `NEXT_PUBLIC_API_URL=http://localhost:1234`.

## Architecture

**API** — Clean Architecture: handlers → services → repositories → models.

| Layer | Location | Role |
|-------|----------|------|
| Routes | `server/routes/routes.go` | Wire endpoints |
| Handlers | `server/handlers/` | Validate requests, call services, return responses |
| Requests/Responses | `server/requests/`, `server/responses/` | camelCase JSON structs |
| Services | `services/{domain}/` | Business logic + `*_test.go` |
| Repositories | `repositories/` | DB access only |
| Models | `db/models/` | GORM |
| Migrations | `db/migrations/` | Numbered up/down pairs — always both |
| Tests | embed `testutils.TestSuite`; fixtures via `s.seeder` |

Every request runs in a DB transaction; services come from context via `h.Services(c)`. Permissions: `services/authorization/` + `middlewares.RequirePermissions()`.

**Frontend** — App Router under `src/app/` (route groups for layouts). Feature components in `src/components/`; shared UI in `src/components/ui/`. Redux Toolkit in `src/store/`; domain axios clients in `src/api/`. Use `useAPI` in components; axios from Redux; thunks for global state flows. Styling: Tailwind + shadcn. Forms: react-hook-form + yup/zod. Unit tests: `src/components/ui/__tests__/` (no snapshots; prefer `user-event`).

## How to add things

- **API endpoint**: mirror a similar handler → request/response → service → route → service tests.
- **DB table**: migration up+down → model → repository → seeder if needed → migrate.
- **Frontend page**: route under `src/app/` → components → `src/api/{domain}.ts` → Redux only if global state is required.

## What to test

Write the **minimum** tests that give confidence shipping. Coverage % does not matter. Prefer one meaningful behavioral or journey test over many tiny implementation tests. When unsure whether a behavior is worth testing, ask.

### Write tests for

- User journeys and core product flows (auth, onboarding, create/edit, invite, upload/download, checkout, password reset)
- Business rules: permissions, calculations, pricing, eligibility, parsers/transformers, state machines, scheduling, retries
- API contracts, validation, serialization, persistence, migrations
- Navigation/routing, forms, search/filter/sort/pagination
- Concurrency-sensitive behavior: races, optimistic updates, cache invalidation
- Destructive actions, undo/recovery, security-sensitive paths
- Feature flags when they change real outcomes
- Interactive accessibility where it matters: keyboard navigation, focus management, dialogs/menus/dropdowns — not decorative markup

### Unit tests

Target **business logic** and observable outcomes. Mock only when isolation is required. Do not test UI copy, labels, CSS, DOM structure, snapshots, icons, React/hook/reducer internals, mock call counts, framework/library behavior, types, constants, or trivial getters.

### E2E tests

Executable user journeys only — prove a user can complete something meaningful end to end. Do not E2E static pages, marketing copy, element presence, styling, or animations.

### Style

Broad behavioral tests over many narrow ones. Shared setup/helpers over duplicated fixtures. Skip edge cases already covered elsewhere.

## Before committing

```bash
cd api && go mod tidy && go build -o build/main cmd/main.go && sh ./scripts/setup_test_db.sh && go test -v ./...
cd frontend && pnpm run type-check && pnpm run lint && pnpm run test
```

For schema changes: test migrate up and down; update `tests/testutils/seeder.go` if fixtures change.
