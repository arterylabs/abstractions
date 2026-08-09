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

| Variable                           | Required | Description                                                                   |
| ---------------------------------- | -------- | ----------------------------------------------------------------------------- |
| `NEXT_PUBLIC_API_URL`              | Yes      | Go API base URL (e.g. `http://localhost:1234`)                                |
| `NEXT_PUBLIC_SITE_URL`             | No       | Hostname without protocol for SEO metadata; defaults to `useabstractions.com` |
| `NEXT_PUBLIC_ANALYTICS_ENABLED`    | No       | Must be exactly `true` to enable PostHog web analytics (off by default)       |
| `NEXT_PUBLIC_ANALYTICS_KEY`        | No       | PostHog project key; required when analytics is enabled                       |
| `NEXT_PUBLIC_ANALYTICS_PROXY_PATH` | No       | First-party proxy path used as SDK `api_host` (default `/p`)                  |

Analytics is **fail closed**: off unless both the enabled flag and a key are set. The SDK uses a first-party proxy path (default `/p`) as `api_host` — proxy that path at the edge if you enable analytics. Self-hosted deployments should leave analytics unset.

## Layout

| Path                | Role                                                   |
| ------------------- | ------------------------------------------------------ |
| `src/app/`          | App Router routes (auth, dashboard, content, god-mode) |
| `src/components/`   | UI and feature components                              |
| `src/api/`          | Axios API clients                                      |
| `src/store/`        | Redux Toolkit                                          |
| `src/content/help/` | In-app help (MDX)                                      |
| `e2e/`              | Playwright specs                                       |
