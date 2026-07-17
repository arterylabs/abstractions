<div align="center">

<!-- Light text for dark theme · dark text for light theme -->
<img src="frontend/public/assets/logo.svg#gh-dark-mode-only" alt="Abstractions" width="280" />
<img src="frontend/public/assets/logo-light.svg#gh-light-mode-only" alt="Abstractions" width="280" />

---

**Weekly AI newsletters that explain your codebase.**

Connect GitHub, index your repos, and keep the whole team — humans and agents — aligned as the code changes.

[Website](https://useabstractions.com) · [Self-host](#quick-start) · [AGPL-3.0](./LICENSE)

</div>

---

## Why

I recently joined a new company and used Claude extensively to build my understanding of the large monolithic codebase. Being me, I was also curious about how other parts of the codebase worked — the parts I don't touch day to day. So I built a tool that doesn't depend on me asking questions to feed that curiosity. Abstractions indexes your codebase and sends out a weekly newsletter to you and your team about how different functionalities or concepts work.

You might already ask Claude to answer your questions about your codebase. Abstractions works the same way, but it doesn't wait for you to ask. It understands the code, figures out a topic, and sends out a newsletter every week — explaining **what changed**, **how the system fits together**, and **why it was built that way**.

|               |                                                                                 |
| ------------- | ------------------------------------------------------------------------------- |
| **Hosted**    | [useabstractions.com](https://useabstractions.com) — same codebase as this repo |
| **Self-host** | Run the stack on your infrastructure under AGPL-3.0                             |
| **Open**      | Inspect indexing, embeddings, and newsletter generation end to end              |

## Features

- **GitHub-native** — OAuth, repo connection, and webhooks for ongoing updates
- **Semantic index** — file chunking + embeddings (Postgres + pgvector)
- **AI newsletters** — topics and weekly digests that stay current with the code
- **Workspaces** — members, invitations, and role-based access
- **Your keys or ours** — Anthropic, OpenAI, or xAI/Grok (BYOK or platform defaults)
- **Delivery that respects people** — subscribers, suppression, and transactional email

## Architecture

Monorepo with two independent apps (no shared runtime code):

| Path        | Stack                              | Role                          |
| ----------- | ---------------------------------- | ----------------------------- |
| `api/`      | Go · Echo · GORM · Asynq           | REST API + background workers |
| `frontend/` | Next.js 15 · React 19 · TypeScript | Web UI (App Router)           |

```
Browser → Next.js → Go API → Postgres (+ pgvector)
                        └── Redis → Workers (index · summarize · newsletter)
```

## Prerequisites

- **Go** (see `api/go.mod` toolchain)
- **Node.js** 22+ and **pnpm** 9.15+
- **PostgreSQL** 15+ with [pgvector](https://github.com/pgvector/pgvector)
- **Redis** 7+
- Optional: Cloudflare Email Service (or AWS SES) + S3, Google/GitHub OAuth apps, AI provider API keys

## Quick start

### 1. Configure the API

```bash
cp api/.env.example api/.env
```

Set at least `SIGNING_KEY`, `TOKEN_REFRESH_KEY`, `ENCRYPTION_KEY_32_BIT` (exactly 32 characters), and your database/Redis settings. Full list: [`api/.env.example`](api/.env.example).

### 2. Database

```bash
# Create the DB, then enable pgvector:
#   CREATE EXTENSION IF NOT EXISTS vector;

cd api
go mod tidy
go run cmd/main.go migrate 0
```

### 3. API + workers

```bash
# Terminal 1 — API (http://localhost:1234)
cd api && go run cmd/main.go

# Terminal 2 — workers (indexing + newsletters)
cd api && go run workers/workers.go
```

### 4. Frontend

```bash
cp frontend/.env.example frontend/.env.local
cd frontend && pnpm install && pnpm run dev --turbopack
```

Open [http://localhost:3000](http://localhost:3000).

<details>
<summary><strong>Docker (optional)</strong></summary>

```bash
# Uses docker-compose.override.yml for local image builds.
# Ensure api/.env exists. Optionally:
#   export NEXT_PUBLIC_API_URL=http://localhost:1234
docker compose up --build
```

See [`.env.docker.example`](.env.docker.example) for host networking notes when the API runs in containers against host Postgres/Redis.

</details>

## Development

```bash
./setup.sh                    # pre-commit: gofmt + ESLint/Prettier
sh api/scripts/gofmt.sh       # format all API Go files
```

**Tests**

```bash
# API (Postgres required)
cd api && sh ./scripts/setup_test_db.sh && go test -v ./...

# Frontend
cd frontend && pnpm run test
```

CI: [`.github/workflows/tests.yml`](.github/workflows/tests.yml).

Agent-oriented conventions live in [`AGENTS.md`](AGENTS.md). In-app help content is under `frontend/src/content/help/`.

## Self-hosting

Point `CLIENT_URL`, `API_URL`, and `NEXT_PUBLIC_SITE_URL` at your domains. Use `ENV=production` on the public internet (disables test-only endpoints and enables auth rate limits).

This repo ships **useabstractions.com** branding. If you run your own instance, plan to rebrand:

| Area                              | Location                                            |
| --------------------------------- | --------------------------------------------------- |
| Email sender, support, legal name | `api/clients/email/email.go`                        |
| Newsletter footer links           | `api/workers/jobs/generate_newsletter.go`           |
| Legal pages                       | `frontend/src/app/(content)/`                       |
| Site URL / metadata               | `frontend/src/app/layout.tsx`, `shared-metadata.ts` |
| Nginx example                     | `conf/useabstractions.com.conf`                     |

Trademarks and hosted product branding are **not** licensed under the AGPL.

<details>
<summary><strong>Operator notes</strong></summary>

- **God mode** — set `users.god_mode = true` for a trusted admin to access `/god-mode`. No public self-elevate endpoint.
- **Test endpoints** — when `ENV` is not `production`, `POST /test/cleanup` and `POST /test/verify_user` are registered for e2e.
- **E2e mail** — addresses ending in `@test.arterylabs.com` never hit the email provider (Cloudflare or SES).

</details>

## Security

Report vulnerabilities privately (GitHub Security Advisories when available, or contact via [useabstractions.com](https://useabstractions.com)). Do not open public issues that include secrets, production tokens, or private customer data.

## License

Licensed under the [GNU Affero General Public License v3.0](./LICENSE). If you modify this program and run it as a network service, you must make the Corresponding Source available to users under AGPL-3.0 ([section 13](https://www.gnu.org/licenses/agpl-3.0.html)).

Hosted branding for [useabstractions.com](https://useabstractions.com) is excluded from that license; see [Self-hosting](#self-hosting).
