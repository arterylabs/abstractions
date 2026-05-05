---
name: review-tests
description: Reviews all existing Abstractions API test cases, identifies which modules lack tests, flags gaps in coverage for changed modules, and produces an actionable improvement plan. Use this when asked to review tests, audit test coverage, or plan test improvements.
---

You are a senior Go engineer performing a thorough audit of the Abstractions API test suite. Your job is to review every existing test file, identify gaps, and produce a clear improvement plan. Do not make any code changes — only analyse and report.

---

## Step 1: Inventory existing tests

Run the following to map every service to its test files and test functions:

```bash
# List services and whether they have test files
for f in api/services/*/; do
  service=$(basename "$f")
  tests=$(ls "$f"*_test.go 2>/dev/null | wc -l | tr -d ' ')
  echo "$service: $tests test file(s)"
done

# List all test function names per file
for f in $(find api/services -name "*_test.go"); do
  echo "=== $f ==="
  grep "^func (s \*" "$f"
done
```

Also inventory non-service tests:
```bash
find api -name "*_test.go" | grep -v "api/services"
```

---

## Step 2: Inventory service public methods

For each service, list its exported functions to compare against test coverage:

```bash
for f in api/services/*/service.go; do
  echo "=== $f ==="
  grep "^func " "$f"
done
```

---

## Step 3: Gap analysis

For each service, assess:

1. **No tests at all** — the service directory has zero `*_test.go` files
2. **Partial coverage** — service has tests but exported methods are untested
3. **Missing error paths** — only happy paths are tested; no tests for invalid input, not-found, permission errors
4. **Missing edge cases** — boundary conditions, pagination edge cases, concurrent operations

### Known gaps (as of last audit)

| Module | Status | Missing coverage |
|--------|--------|-----------------|
| `services/twilio` | ❌ No tests | `SendSMS`, `ValidateWebhook`, `UsesStatusCallback` |
| `services/whatsapp` | ❌ No tests | `SendMessage`, `DownloadMedia`, `ValidateWebhook` |
| `services/event` | ⚠️ Partial | `GetEvent`, `GetEventBySlug`, `UpdateUserRole`, `RemoveUserFromEvent` |
| `services/guest` | ⚠️ Partial | `ImportGuestsFromPayload`, `GetGuestDetail`, `GetStats`, `SendInvitation` |
| `services/invitation` | ⚠️ Partial | Error paths for `AcceptInvitation` (already-accepted, wrong token) |
| `services/user` | ⚠️ Partial | `UpdateUser`, `GetAdminUsers` list pagination, god-mode flag toggling |
| `services/session` | ⚠️ Partial | Concurrent session creation, refresh token rotation |
| `repositories` | ❌ No integration tests | All repository methods untested (only enum validation exists) |

---

## Step 4: Prioritisation

Rank improvements by impact:

### Priority 1 — Missing test files (high risk, zero coverage)
- `services/twilio` — used for SMS delivery; `ValidateWebhook` is a security-critical signature check
- `services/whatsapp` — same risk profile as Twilio

### Priority 2 — Security & correctness gaps
- `services/authorization` — missing tests for `viewer`/`editor` roles (only `host` and `owner` tested)
- `services/invitation` — no test for expired invitation or duplicate acceptance
- `services/session` — no test for refresh token rotation or concurrent invalidation

### Priority 3 — Business logic completeness
- `services/event` — `GetEventBySlug` (public-facing), `UpdateUserRole`, `RemoveUserFromEvent`
- `services/guest` — `ImportGuestsFromPayload` (bulk import error handling), `GetStats`, `SendInvitation`
- `services/user` — god-mode toggle, admin list pagination

### Priority 4 — Repository layer
- All repository files in `api/repositories/` lack integration tests
- Worth adding for complex query methods (e.g. `GuestRepository.List`, `EventRepository.FindBySlug`)

---

## Step 5: Output format

Produce a structured report:

```
## Test Coverage Audit — Abstractions API

### Summary
- X service packages: Y fully covered, Z partial, W no tests
- X utility packages: all covered
- Repositories: no integration tests

### Services with no tests
- services/twilio
- services/whatsapp

### Services with partial coverage
For each: list tested methods and untested methods

### Recommendations
Ordered list of test additions with rationale
```

---

## Step 6: Do NOT make changes

This skill is read-only. Output the report and recommendations only. To implement improvements, use the `improve-tests` skill.
