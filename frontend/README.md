# Frontend

Next.js 15 (App Router) UI for Abstractions.

See the [root README](../README.md) for product overview, full local setup, environment variables, and self-hosting notes.

**License:** [GNU Affero General Public License v3.0](../LICENSE) (AGPL-3.0) — same as the monorepo root.

## Quick commands

```bash
# From this directory
cp .env.example .env.local   # then set NEXT_PUBLIC_API_URL

pnpm install
pnpm run dev --turbopack     # http://localhost:3000

pnpm run type-check
pnpm run lint
pnpm run test
pnpm run test:e2e            # requires API + frontend running
```

## Environment

| Variable               | Required | Description                                                                   |
| ---------------------- | -------- | ----------------------------------------------------------------------------- |
| `NEXT_PUBLIC_API_URL`  | Yes      | Go API base URL (e.g. `http://localhost:1234`)                                |
| `NEXT_PUBLIC_SITE_URL` | No       | Hostname without protocol for SEO metadata; defaults to `useabstractions.com` |

## Layout

| Path                | Role                                                   |
| ------------------- | ------------------------------------------------------ |
| `src/app/`          | App Router routes (auth, dashboard, content, god-mode) |
| `src/components/`   | UI and feature components                              |
| `src/api/`          | Axios API clients                                      |
| `src/store/`        | Redux Toolkit                                          |
| `src/content/help/` | In-app help (MDX)                                      |
| `e2e/`              | Playwright specs                                       |
