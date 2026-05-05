---
name: improve-tests
description: Runs the full Abstractions API test suite, fixes any failing tests, then implements new test cases for modules that lack coverage or need improvement. Always get explicit user approval before writing new tests. Use this when asked to fix tests, add missing tests, or improve the test suite.
---

You are a senior Go engineer improving the Abstractions API test suite. You will run the existing tests, fix failures, then propose and (with approval) implement new tests based on coverage gaps.

---

## Phase 1: Run existing tests

### 1.1 — Ensure the test database is ready

```bash
cd api
sh ./scripts/setup_test_db.sh
```

### 1.2 — Run the full suite

```bash
cd api
go test -v ./... 2>&1 | tee /tmp/harmony_test_output.txt
```

### 1.3 — Parse results

```bash
# Summary
grep -E "^(ok|FAIL|---)" /tmp/harmony_test_output.txt

# Only failures
grep -E "^(FAIL|--- FAIL)" /tmp/harmony_test_output.txt
```

---

## Phase 2: Fix failing tests

For each failing test:

1. Read the full failure output and the test file
2. Read the service/code under test to understand the expected behaviour
3. Determine whether the failure is in:
   - **The test itself** (wrong expectation, stale fixture, bad setup)
   - **The service** (genuine bug introduced by a recent change)
4. Fix the test (preferred) or fix the service (only if there is a real bug)
5. Re-run the specific package to confirm the fix:

```bash
cd api && go test -v ./services/<package>
```

Repeat until `go test -v ./...` passes with no `FAIL` lines.

---

## Phase 3: Coverage gap analysis

After all tests pass, audit coverage using the `review-tests` skill patterns:

```bash
# Services with no test files
for f in api/services/*/; do
  service=$(basename "$f")
  tests=$(ls "$f"*_test.go 2>/dev/null | wc -l | tr -d ' ')
  [ "$tests" -eq 0 ] && echo "NO TESTS: $service"
done

# Public methods per service vs test functions
for svc in api/services/*/service.go; do
  echo "=== $(dirname $svc | xargs basename) ==="
  grep "^func (s \*Service)" "$svc" | awk '{print $2}' | sed 's/(.*//'
done
```

---

## Phase 4: Propose improvements

**Before writing any new tests, present the following to the user for approval.**

Format the proposal as a table:

| # | Module | Test name | What it covers | Priority |
|---|--------|-----------|----------------|----------|
| 1 | services/twilio | `TestSendSMS_Success` | Happy path SMS delivery using a mock HTTP server | High |
| 2 | services/twilio | `TestValidateWebhook_ValidSignature` | Correct HMAC-SHA1 signature is accepted | High |
| 3 | services/twilio | `TestValidateWebhook_InvalidSignature` | Tampered signature is rejected | High |
| 4 | services/whatsapp | `TestSendMessage_Success` | Happy path WhatsApp message via mock HTTP | High |
| 5 | services/whatsapp | `TestValidateWebhook_ValidSignature` | Webhook HMAC validation passes | High |
| 6 | services/event | `TestGetEvent_NotFound` | Returns error for non-existent event ID | Medium |
| 7 | services/event | `TestGetEventBySlug_ReturnsCorrectEvent` | Slug lookup returns right event | Medium |
| 8 | services/event | `TestUpdateUserRole_ChangesRole` | Role changes are persisted | Medium |
| 9 | services/event | `TestRemoveUserFromEvent_RemovesRelationship` | User is removed, event still exists | Medium |
| 10 | services/guest | `TestImportGuestsFromPayload_DeduplicatesPhones` | Bulk import skips duplicates | Medium |
| 11 | services/guest | `TestGetStats_CountsRSVPStatuses` | Stats match created guests | Medium |
| 12 | services/invitation | `TestAcceptInvitation_AlreadyAccepted` | Second accept returns error | Medium |
| 13 | services/invitation | `TestAcceptInvitation_ExpiredToken` | Expired invitation is rejected | Medium |
| 14 | services/authorization | `TestEditorPermissions` | Editor role has correct permission set | Medium |
| 15 | services/authorization | `TestViewerPermissions` | Viewer role has correct permission set | Medium |
| 16 | services/session | `TestSessionInvalidation_Concurrent` | Invalidating used session blocks reuse | Low |

Then ask: **"Would you like me to implement all of these, or select specific ones?"**

---

## Phase 5: Implement approved tests

Only proceed after explicit user approval.

For each approved test:

### Writing Twilio/WhatsApp tests (mock HTTP)

These services make outbound HTTP calls. Use `net/http/httptest` to mock Twilio's API:

```go
package twilio_service_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    twilio_service "abstractions/services/twilio"
    "github.com/stretchr/testify/assert"
)

func TestSendSMS_Success(t *testing.T) {
    // Mock Twilio API server
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"sid":"SM123","status":"queued"}`))
    }))
    defer srv.Close()

    // Inject mock base URL via exported setter or constructor option
    svc := twilio_service.NewServiceWithBaseURL(srv.URL, "ACtest", "authtest", "+15550000000", "", "")
    result, err := svc.SendSMS(context.Background(), "+15551234567", "Hello!")
    assert.NoError(t, err)
    assert.Equal(t, "SM123", result.MessageSID)
}
```

> **Note**: If `NewServiceWithBaseURL` does not exist, add it to `api/services/twilio/service.go` as part of this improvement — a constructor that accepts a `baseURL` override enables testing without external network calls.

### Writing service tests that need the DB

Embed `testutils.TestSuite` as normal — see the `test-setup` skill for the full pattern.

### File placement

| Service | Test file path |
|---------|---------------|
| `services/twilio` | `api/services/twilio/twilio_service_test.go` |
| `services/whatsapp` | `api/services/whatsapp/whatsapp_service_test.go` |
| `services/event` | Add methods to `api/services/event/event_service_test.go` |
| `services/guest` | Add methods to `api/services/guest/guest_service_test.go` |
| `services/invitation` | Add methods to `api/services/invitation/invitation_service_test.go` |
| `services/authorization` | Add methods to `api/services/authorization/authorization_service_test.go` |

---

## Phase 6: Verify

After implementing each batch of tests:

```bash
cd api && go test -v ./services/<package>
```

Then run the full suite to confirm nothing regressed:

```bash
cd api && go test -v ./...
```

Only mark the task complete when `go test -v ./...` exits with code 0.

---

## Rules

- **Always get approval before writing new tests** (Phase 4 gate)
- Never modify production service code solely to make a test pass — only fix genuine bugs
- If a service lacks a testable constructor (e.g. no way to inject a mock HTTP client), add one as a minimal, non-breaking change and note it in the proposal
- Use `assert.*` (non-fatal) unless a nil pointer would cause a panic on the next line; use `s.Require()` sparingly
- All new DB-touching tests must use `s.DB` and `s.Ctx` from `testutils.TestSuite`
