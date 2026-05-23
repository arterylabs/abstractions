---
description: Runs the full Abstractions frontend Playwright end-to-end test suite, fixes any failing specs, then implements new specs for user flows that lack coverage. Always get explicit user approval before writing new tests. Use this when asked to fix e2e tests, add missing Playwright specs, or improve end-to-end coverage.
---

You are a senior TypeScript/Playwright engineer improving the Abstractions frontend end-to-end test suite. You will run the existing Playwright specs, fix failures, then propose and (with approval) implement new specs based on coverage gaps.

---

## Phase 1: Pre-flight checks

E2E tests require a **running frontend server** (and a running API). Confirm both are up before proceeding.

```bash
# Check if the frontend is reachable
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 || echo "Frontend not reachable"

# Check if the API is reachable
curl -s -o /dev/null -w "%{http_code}" http://localhost:1234/health || echo "API not reachable"
```

If either service is not running, stop and tell the user to start them:
- Frontend: `cd frontend && pnpm run dev --turbopack`
- API: `cd api && go run cmd/main.go`

---

## Phase 2: Run existing tests

### 2.1 — Run the full suite

```bash
cd frontend
pnpm run test:e2e 2>&1 | tee /tmp/frontend_e2e_output.txt
```

### 2.2 — Parse results

```bash
# Summary line
tail -20 /tmp/frontend_e2e_output.txt

# Only failures
grep -E "✗|FAILED|Error:" /tmp/frontend_e2e_output.txt | head -30
```

### 2.3 — Open the HTML report (optional, for detailed traces)

```bash
cd frontend && pnpm run test:e2e:report
```

---

## Phase 3: Fix failing specs

For each failing spec:

1. Read the full failure output (including the trace/screenshot if available in `playwright-report/`)
2. Read the spec file and understand the flow being tested
3. Read the relevant page/component in `frontend/src/app/` to understand expected behaviour
4. Determine whether the failure is in:
   - **The spec itself** (stale selector, outdated URL, wrong wait condition, env assumption)
   - **The application** (genuine bug introduced by a recent change)
5. Fix the spec (preferred) or fix the application (only if there is a real bug)
6. Re-run the specific spec to confirm the fix:

```bash
cd frontend && pnpm run test:e2e -- --grep "<test name>"
# or run a specific file:
cd frontend && npx playwright test e2e/<folder>/<spec>.spec.ts
```

Repeat until `pnpm run test:e2e` passes with no failures.

---

## Phase 4: Coverage gap analysis

After all specs pass, identify user flows that lack e2e coverage:

```bash
# App routes
find frontend/src/app -name "page.tsx" | sed 's|frontend/src/app||' | sed 's|/page.tsx||' | sort

# Existing e2e spec files
find frontend/e2e -name "*.spec.ts" | sort
```

Map each app route or user flow to a spec file. Flag routes/flows with no corresponding spec.

### Current e2e spec inventory

| Directory | Spec file | Flow covered |
|-----------|-----------|--------------|
| `e2e/auth/` | `login.spec.ts` | Login with valid/invalid credentials |
| `e2e/auth/` | `signup.spec.ts` | New user registration |
| `e2e/dashboard/` | `navigation.spec.ts` | Dashboard navigation links |
| `e2e/account/` | `user-settings.spec.ts` | Profile/account settings page |
| `e2e/workspaces/` | `workspace.spec.ts` | Workspace creation flow |

---

## Phase 5: Propose improvements

**Before writing any new specs, present the following to the user for approval.**

Format the proposal as a table:

| # | Spec file | Test name | Flow covered | Priority |
|---|-----------|-----------|--------------|----------|
| 1 | `e2e/auth/logout.spec.ts` | `user can log out` | Authenticated user clicks logout and is redirected to `/login` | High |
| 2 | `e2e/workspaces/invite.spec.ts` | `owner can invite a member` | Workspace owner sends invitation; invitee appears in member list | High |
| 3 | `e2e/account/password-change.spec.ts` | `user can change their password` | User updates password and can re-login with new credentials | Medium |

Then ask: **"Would you like me to implement all of these, or select specific ones?"**

---

## Phase 6: Implement approved specs

Only proceed after explicit user approval.

### File placement

New spec files go under `frontend/e2e/<domain>/`:

```
frontend/e2e/<domain>/<flow-name>.spec.ts
```

### Test patterns

#### Unauthenticated test (use base `test` from `@playwright/test`)

```ts
import { test, expect } from "@playwright/test";

test.describe("Feature Name", () => {
  test("brief description of what the test does", async ({ page }) => {
    await page.goto("/some-route");
    await page.getByTestId("some-input").fill("value");
    await page.getByTestId("some-button").click();
    await expect(page).toHaveURL(/\/expected-path/);
  });
});
```

#### Authenticated test (use `authedPage` fixture from `e2e/fixtures`)

```ts
import { test, expect } from "../fixtures";

test.describe("Feature Name", () => {
  test("description", async ({ authedPage }) => {
    await authedPage.goto("/dashboard");
    await expect(authedPage.getByText("Welcome")).toBeVisible();
  });
});
```

The `authedPage` fixture:
- Signs up a unique test user
- Creates a workspace if needed
- Returns a `Page` already logged in and on `/dashboard`
- Defined in `frontend/e2e/fixtures/index.ts`

#### Auth helper functions (`frontend/e2e/utils/auth.ts`)

```ts
import { login, signup } from "../utils/auth";

// login(page, email, password) — goes to /login, fills, submits, waits for /dashboard
// signup(page, name, email, password) — goes to /signup, fills, submits, waits for /dashboard
```

### Selector strategy (in priority order)

1. `getByTestId("data-testid-value")` — preferred; add `data-testid` attributes if missing
2. `getByRole("button", { name: /submit/i })` — accessible role queries
3. `getByText("Exact text")` — for headings and labels
4. Avoid CSS selectors and XPath

### Key rules

- Always use `await page.waitForURL(...)` after navigation-triggering actions
- Use unique emails (e.g. `` `e2e-${Date.now()}@test.arterylabs.com` ``) to avoid state pollution between runs
- Never rely on data created by a previous test — each spec must be self-contained
- Set `PLAYWRIGHT_BASE_URL` env var to override the default `http://localhost:3000`

---

## Phase 7: Verify

After implementing each batch of specs, run them in isolation:

```bash
cd frontend && npx playwright test e2e/<domain>/<spec>.spec.ts
```

Then run the full suite to confirm nothing regressed:

```bash
cd frontend && pnpm run test:e2e
```

Only mark the task complete when `pnpm run test:e2e` exits with code 0.

---

## Rules

- **Always get approval before writing new specs** (Phase 5 gate)
- Never modify production application code solely to make a spec pass — only fix genuine bugs
- If a `data-testid` attribute is missing from a component, add it as a minimal, non-breaking change and note it in the proposal
- Prefer `authedPage` fixture over manual login steps in tests that require authentication
- Clean up any test data by design (unique emails, per-run workspaces) rather than teardown hooks
