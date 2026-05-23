---
name: security-audit
description: Performs a rigorous, adversarial security audit of the Abstractions codebase. Identifies exploitable vulnerabilities and logic flaws across injection threats, authentication, access control, data exposure, input validation, file uploads, and business logic. Use this when asked to audit security, find vulnerabilities, or review the codebase for exploits.
---

You are performing a rigorous, adversarial security audit of the Abstractions codebase. Focus purely on identifying exploitable vulnerabilities and logic flaws — not code style or generic quality improvements.

For each vulnerability found, provide:

- **Vulnerability Type & Location**: File, line number, or function name
- **Description**: Exactly what the flaw is
- **Exploit Scenario**: How an attacker could abuse this flaw
- **Severity**: Critical / High / Medium / Low
- **Remediation**: An exact, secure code snippet to fix the issue

---

## Step 1: Gather codebase context

Orient yourself before diving into specific vulnerability classes.

```bash
# Map the API surface — all registered routes
grep -r "\.GET\|\.POST\|\.PUT\|\.PATCH\|\.DELETE" api/server/routes/ --include="*.go" -n

# List all middleware applied globally vs per-route
cat api/server/routes/routes.go

# Find all handler files
ls api/server/handlers/

# Find all service files
ls api/services/*/service.go 2>/dev/null

# Find all repository files
ls api/repositories/
```

---

## Step 2: Injection Threats

### SQL / ORM injection

```bash
# Look for raw SQL strings or Exec calls with user input
grep -rn "Raw\|Exec\|\.Where(" api/ --include="*.go" | grep -v "_test.go"

# Look for string concatenation inside queries
grep -rn 'fmt\.Sprintf.*WHERE\|fmt\.Sprintf.*SELECT\|fmt\.Sprintf.*INSERT\|fmt\.Sprintf.*UPDATE\|fmt\.Sprintf.*DELETE' api/ --include="*.go"
```

Flag: any `fmt.Sprintf` or string concat building SQL clauses, or `.Raw()` calls with unsanitized values.

### Command injection

```bash
# Find exec.Command or os.system calls
grep -rn "exec\.Command\|syscall\.Exec\|os\.StartProcess" api/ --include="*.go"
```

Flag: any user-controlled input flowing into command arguments.

### NoSQL / Redis injection

```bash
# Find Redis key construction
grep -rn "\.Set\|\.Get\|\.Del\|\.HSet\|\.HGet" api/ --include="*.go" | grep -v "_test.go"
```

Flag: user-controlled strings directly used as Redis keys without sanitization or namespacing.

---

## Step 3: Authentication & Session Management

```bash
# Find password hashing — must use bcrypt/argon2/scrypt
grep -rn "bcrypt\|argon2\|scrypt\|sha256\|sha1\|md5\|base64" api/ --include="*.go" | grep -iv "_test.go"

# Find JWT creation and signing
grep -rn "jwt\.\|SigningMethod\|NewWithClaims\|ParseWithClaims" api/ --include="*.go"

# Check token expiration
grep -rn "ExpiresAt\|exp\b" api/ --include="*.go"

# Find login rate-limiting
grep -rn "RateLimit\|limiter\|throttle" api/ --include="*.go"

# Find session creation
grep -rn "CreateSession\|NewSession\|session\.Create" api/ --include="*.go"
```

Flag:
- Passwords stored as plain text or hashed with MD5/SHA1/SHA256 (must be bcrypt, argon2, or scrypt)
- JWTs without expiration (`ExpiresAt` unset or far future)
- JWTs signed with `alg: none` or weak HMAC secrets
- No rate limiting on `/login`, `/signup`, `/verify`, `/reset-password`
- Refresh tokens that never expire or aren't rotated on use

---

## Step 4: Broken Access Control

```bash
# Find all routes and their middleware chains
cat api/server/routes/routes.go

# Find authorization checks in handlers
grep -rn "RequirePermissions\|AuthMiddleware\|IsAdmin\|GodMode\|Forbidden\|Unauthorized" api/server/ --include="*.go" -n

# Find places where user IDs are taken directly from request params (IDOR risk)
grep -rn "Param(\"id\")\|Param(\"userId\")\|Param(\"workspaceId\")" api/server/handlers/ --include="*.go" -n

# Find ownership checks — does the service verify the requester owns the resource?
grep -rn "UserID\|WorkspaceID\|OwnedBy\|CreatedBy" api/services/ --include="*.go" -n
```

Flag:
- Routes without `AuthMiddleware` that should be protected
- Handlers that fetch a resource by ID from the URL without checking the caller owns it (IDOR)
- Any route where a regular user can pass `isAdmin=true` or `godMode=true` in a request body
- Privilege escalation paths (e.g., a user updating their own role to `owner`)

---

## Step 5: Data Exposure

```bash
# Find log statements that may print sensitive fields
grep -rn "log\.\|logger\.\|fmt\.Print\|zerolog" api/ --include="*.go" | grep -i "password\|token\|secret\|key\|hash\|email" | grep -v "_test.go"

# Find response structs — do they expose internal fields?
cat api/server/responses/*.go

# Find error responses that return raw DB or internal errors
grep -rn "c\.JSON.*err\|c\.JSON.*error\b" api/server/handlers/ --include="*.go" -n

# Check .env.example for hardcoded secrets
cat api/.env.example 2>/dev/null || cat .env.docker.example 2>/dev/null
```

Flag:
- Passwords, tokens, or private keys included in JSON responses
- Raw database error messages returned to the client (leaks schema info)
- Logging of sensitive fields (passwords, tokens)
- Environment variables or secrets committed to source files

---

## Step 6: Input Validation

```bash
# Find request validation structs
cat api/server/requests/*.go

# Find where request bodies are bound — is Validate() called?
grep -rn "c\.Bind\|c\.Validate\|ShouldBind" api/server/handlers/ --include="*.go" -n

# Look for missing length limits on text fields
grep -rn "validate:\"" api/server/requests/ --include="*.go"

# Find frontend forms — does the API enforce size limits?
grep -rn "MaxLength\|maxlength\|max=" frontend/src/ --include="*.tsx" --include="*.ts" -n | head -30
```

Flag:
- Request structs missing `validate` tags for required fields, max length, or format
- Handlers that call `c.Bind()` without `c.Validate()`
- No maximum payload size configured on the HTTP server
- User-controlled strings rendered in the frontend without escaping (XSS)

---

## Step 7: File Upload Vulnerabilities

```bash
# Find file upload handlers
grep -rn "FormFile\|MultipartForm\|Upload\|storage\." api/ --include="*.go" -n

# Check MIME type validation
grep -rn "ContentType\|MimeType\|DetectContentType\|mime\." api/ --include="*.go" -n

# Check for path sanitization
grep -rn "filepath\.Join\|path\.Join\|Filename" api/ --include="*.go" -n
```

Flag:
- Missing MIME type validation (only checking file extension is insufficient)
- `filepath.Join` with user-supplied filename without `filepath.Clean` and base-dir check (path traversal)
- No restriction on executable file types (`.exe`, `.sh`, `.php`, `.py`)
- No file size limit enforced before reading the upload into memory

---

## Step 8: Business Logic Flaws

```bash
# Find payment / subscription flows
grep -rn "payment\|stripe\|charge\|subscribe\|invoice\|billing" api/ --include="*.go" -i -n

# Find idempotency checks
grep -rn "idempotent\|duplicate\|already.*exists\|conflict" api/services/ --include="*.go" -i -n

# Find invitation / token flows
grep -rn "InvitationToken\|VerifyToken\|AcceptInvitation\|UseToken" api/ --include="*.go" -n

# Find rate limiting in business flows
grep -rn "RateLimit\|limiter\|cooldown\|throttle" api/ --include="*.go" -n
```

Flag:
- Invitation or verification tokens that can be used multiple times (no single-use enforcement)
- Lack of idempotency on payment or subscription endpoints
- Workflow steps that can be skipped by calling later-stage endpoints directly
- Email verification that can be bypassed via direct API calls

---

## Step 9: AI / Agentic Risks (if applicable)

```bash
# Find AI/LLM integration points
grep -rn "openai\|anthropic\|llm\|prompt\|completion\|embedding" api/ --include="*.go" -i -n
grep -rn "openai\|anthropic\|llm\|prompt\|completion\|embedding" frontend/src/ --include="*.ts" --include="*.tsx" -i -n

# Find places where external data is used in prompts
grep -rn "UserMessage\|SystemPrompt\|context\[" api/ --include="*.go" -n
```

Flag:
- User-controlled text injected directly into system prompts (prompt injection)
- Tool execution triggered by model output without explicit allow-list validation
- Retrieval-augmented data inserted into prompts without sanitization
- No output filtering before displaying model responses to users

---

## Step 10: Output format

Produce a structured report:

```
## Security Audit Report — Abstractions

### Executive Summary
- Total vulnerabilities: X (Critical: N, High: N, Medium: N, Low: N)
- Highest-risk areas: [list areas]

### Findings

#### [SEVERITY] Vulnerability Title
- **Location**: file:line or function
- **Description**: ...
- **Exploit Scenario**: ...
- **Remediation**:
  ```go
  // secure replacement code
  ```

### No Issues Found In
- [List areas that were clean]

### Recommendations (Priority Order)
1. ...
2. ...
```

---

## Important constraints

- **Do NOT make any code changes** during this audit. Output findings and remediation snippets only.
- If a vulnerability has already been remediated, note it as "Fixed" and skip.
- Prioritise findings by exploitability, not just theoretical severity.
- If a finding requires more context to confirm, note it as "Needs verification" with the specific question.
