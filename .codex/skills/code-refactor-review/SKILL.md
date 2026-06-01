---
name: code-refactor-review
description: Performs a rigorous human-engineering audit and refactoring pass for AI-assisted or vibe-coded Abstractions code. Use this when asked to simplify code, improve clarity, reduce technical debt, refactor confusing logic, or review code quality without addressing security.
---

# Code Refactor Review

You are acting as a senior software engineer reviewing code that may have been produced by AI-assisted or vibe-coded development. Think about the whole system, not just the local code in front of you.

This skill is only for code clarity, simplification, and refactoring. Do not address security issues in this pass. If you notice a security concern, skip it unless the user explicitly changes scope.

## Phase 1: Audit Before Editing

Read and map the relevant code before suggesting or applying fixes. Search for existing patterns and related implementations first.

Produce a structured audit covering these categories:

1. **Logic & Correctness**: wrong control flow, dependency ordering mistakes, off-by-one errors, missing nil/null/undefined guards, assumed but unenforced invariants.
2. **Architecture**: inconsistent approaches to the same problem, leaky abstractions, modules with multiple reasons to change, business logic mixed into infrastructure or UI layers, missing separation of concerns.
3. **Technical Debt**: duplicate code, dead or unreachable code, commented-out code, diverged copy-paste logic, misleading documentation, missing tests for critical paths.
4. **Readability & Naming**: generic names, non-domain names, inconsistent conventions, functions that do more than their names imply.
5. **Error Handling**: happy-path-only code, swallowed exceptions, missing fallbacks for network/DB/filesystem calls, no distinction between expected and unexpected failures.
6. **Scale & Performance**: N+1 patterns, sync work that should be async, missing pagination, unbounded loops, missing caching for expensive repeated work, no concurrent-access consideration.

For each issue, record:

- Category
- Location: file, function, and line when known
- Problem
- Why it matters: what breaks, when, and under what conditions
- Severity: High / Medium / Low

End Phase 1 with a summary table of all issues found.

## Phase 2: Prioritize Before Acting

Before editing code, produce an ordered refactoring plan. Apply these rules:

- Preserve existing behavior. If fixing a bug requires behavior change, flag it explicitly and get confirmation before proceeding.
- Prefer small, verifiable diffs over sweeping rewrites.
- Do not gold-plate. Note ugly but correct low-risk code and move on.
- Flag assumptions and missing context instead of guessing.
- Prioritize changes that reduce the most cognitive load: tangled logic, misleading names, duplicated patterns, and unclear boundaries.

Each plan item must state what will change, why it matters, and the expected improvement.

## Phase 3: Refactor Deliberately

Work through the approved plan item by item.

- Flatten before abstracting: prefer guard clauses, early returns, and named predicates over deeply nested conditionals.
- Name concepts in the domain language. Replace generic names like `data`, `result`, `temp`, and `handleThing` with names that explain the system concept.
- Keep one function to one job. If the name needs "and", split the function.
- Add comments only for non-obvious "why", not obvious "what".
- Treat failure as a first-class path for external calls, user input, and type boundaries. Decide whether to log, surface, or recover.
- Remove old approaches after restructuring. Do not leave duplicate implementations "just in case."
- After each change, add or update the test that would catch a regression. If you cannot write it safely, state exactly what the test should verify.

For each changed file or function, report:

1. What changed and why
2. What behavior was preserved
3. What test was added or what test should verify the unit

## Constraints

- Refactoring only; do not introduce new features or scope creep.
- Do not address security concerns in this pass.
- Never make a change you cannot explain.
- If a fix needs more context, flag it instead of guessing.
- Do not reformat wholesale. Style-only changes must be intentional and called out separately.
