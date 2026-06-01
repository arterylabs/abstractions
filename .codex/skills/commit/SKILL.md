---
name: commit
description: Generates a conventional commit message from staged and unstaged changes, stages all changes with git add, and commits them. Use this when asked to commit changes, stage and commit, write a commit message, or create a git commit.
---

When asked to commit changes, follow these steps precisely:

## Step 1: Gather context

Run the following commands to understand the full scope of changes:

```bash
git status
git diff HEAD
git diff --cached
```

Also check for untracked files:

```bash
git ls-files --others --exclude-standard
```

## Step 2: Analyze the changes

Review the diff output carefully to understand:
- Which files were added, modified, or deleted
- The nature of the changes (feature, fix, refactor, chore, etc.)
- The scope/area of the codebase affected (e.g., `api`, `frontend`, `auth`, `db`)

## Step 3: Generate the commit message

Construct a commit message following the **Conventional Commits** specification:

### Format

```
<type>(<scope>): <short summary>

<body>

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

### Rules

- **type** must be one of:
  - `feat` – a new feature
  - `fix` – a bug fix
  - `refactor` – code change that neither fixes a bug nor adds a feature
  - `chore` – build process, tooling, dependency updates
  - `docs` – documentation only changes
  - `style` – formatting, missing semi-colons, etc. (no logic change)
  - `test` – adding or correcting tests
  - `perf` – performance improvement
  - `ci` – CI/CD configuration changes
  - `build` – build system changes
  - `revert` – reverts a previous commit

- **scope** (optional but preferred): the module, package, or area affected — e.g., `auth`, `api`, `frontend`, `db`, `migrations`, `handlers`, `services`

- **short summary**: imperative mood, lowercase, no period, ≤72 characters

- **body** (required): explain *what* changed and *why*. Use bullet points if multiple things changed. Be specific — reference file names, function names, or concepts where helpful. Keep each line ≤72 characters.

- Always append the Co-authored-by trailer on its own line at the end.

- Use `BREAKING CHANGE:` in the body footer if the change is breaking.

### Examples

```
feat(auth): add Google OAuth login support

- Add GoogleOAuthHandler in server/handlers/auth_handler.go
- Register /auth/google and /auth/google/callback routes
- Store OAuth tokens in sessions table
- Add GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET env vars

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

```
fix(api): return 404 instead of 500 when workspace not found

WorkspaceHandler.Get was calling log.Fatal on a not-found error,
crashing the server. Changed to return a proper 404 JSON response
using the existing NotFoundError helper.

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

```
chore(deps): update Go dependencies and run go mod tidy

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

## Step 4: Stage all changes

```bash
git add -A
```

## Step 5: Commit

Use `git commit` with the `-m` flag. Pass the full multi-line message using `$'...'` syntax to preserve newlines:

```bash
git commit -m $'<type>(<scope>): <summary>\n\n<body>\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'
```

## Step 6: Confirm

Show the result of the commit:

```bash
git --no-pager log -1 --stat
```

Report the commit hash and summary to the user.
