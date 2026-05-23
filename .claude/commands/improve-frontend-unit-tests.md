---
description: Runs the full Abstractions frontend Jest unit test suite, fixes any failing tests, then implements new test cases for UI components that lack coverage or need improvement. Always get explicit user approval before writing new tests. Use this when asked to fix frontend unit tests, add missing component tests, or improve the Jest test suite.
---

You are a senior TypeScript/React engineer improving the Abstractions frontend unit test suite. You will run the existing Jest tests, fix failures, then propose and (with approval) implement new tests based on coverage gaps.

---

## Phase 1: Run existing tests

### 1.1 — Run the full suite with coverage

```bash
cd frontend
pnpm run test:coverage 2>&1 | tee /tmp/frontend_unit_test_output.txt
```

### 1.2 — Parse results

```bash
# Summary: pass/fail counts
grep -E "Tests:|Test Suites:|Snapshots:" /tmp/frontend_unit_test_output.txt

# Only failing test suites
grep -E "FAIL " /tmp/frontend_unit_test_output.txt

# Coverage summary table
grep -A 999 "Coverage summary" /tmp/frontend_unit_test_output.txt | head -30
```

---

## Phase 2: Fix failing tests

For each failing test:

1. Read the full failure output and the test file
2. Read the component under test to understand the expected behaviour
3. Determine whether the failure is in:
   - **The test itself** (wrong expectation, stale snapshot, outdated mock, missing jest.mock)
   - **The component** (genuine bug introduced by a recent change)
4. Fix the test (preferred) or fix the component (only if there is a real bug)
5. Re-run the specific test file to confirm the fix:

```bash
cd frontend && pnpm run test -- --testPathPattern="<component-name>.test.tsx"
```

Repeat until `pnpm run test` passes with no `FAIL` lines.

---

## Phase 3: Coverage gap analysis

After all tests pass, identify UI components that lack test files:

```bash
# Components in ui/ that have no corresponding test file
for f in frontend/src/components/ui/*.tsx; do
  name=$(basename "$f" .tsx)
  test_file="frontend/src/components/ui/__tests__/${name}.test.tsx"
  [ ! -f "$test_file" ] && echo "NO TEST: $name"
done
```

Also check coverage thresholds (branches, functions, lines, statements must all be ≥ 60%):

```bash
cd frontend && pnpm run test:coverage -- --coverageReporters=text-summary 2>&1 | tail -20
```

---

## Phase 4: Propose improvements

**Before writing any new tests, present the following to the user for approval.**

Format the proposal as a table:

| # | Component | Test name | What it covers | Priority |
|---|-----------|-----------|----------------|----------|
| 1 | `ui/color-picker` | `TestColorPicker_SelectsColor` | Clicking a swatch fires `onChange` with the correct hex value | High |
| 2 | `ui/data-table` | `TestDataTable_SortsByColumn` | Clicking a sortable column header re-orders rows | High |
| 3 | `ui/datetime-picker` | `TestDatetimePicker_SelectsDate` | Selecting a date from the calendar updates the displayed value | Medium |

Then ask: **"Would you like me to implement all of these, or select specific ones?"**

---

## Phase 5: Implement approved tests

Only proceed after explicit user approval.

### File placement

All new unit tests go in:

```
frontend/src/components/ui/__tests__/<component-name>.test.tsx
```

### Test pattern

Follow the existing convention (no snapshot tests; prefer `@testing-library/react` queries and `@testing-library/user-event` for interactions):

```tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MyComponent } from "../my-component";

describe("MyComponent", () => {
  it("renders without crashing", () => {
    render(<MyComponent />);
    expect(screen.getByRole("...")).toBeInTheDocument();
  });

  it("calls onChange when interacted with", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<MyComponent onChange={onChange} />);
    await user.click(screen.getByRole("button"));
    expect(onChange).toHaveBeenCalledTimes(1);
  });
});
```

**Key rules:**
- Use `screen.getByRole` / `getByText` / `getByTestId` — avoid `container.querySelector`
- Use `userEvent.setup()` (not `fireEvent`) for all user interactions
- Mock external dependencies (e.g. `next/navigation`, `next/image`) at the top of the file with `jest.mock(...)`
- Do not use snapshot tests
- Use `@/` alias for imports within `src/` (e.g. `import { cn } from "@/lib/utils"`)

---

## Phase 6: Verify

After implementing each batch of tests, run the targeted file:

```bash
cd frontend && pnpm run test -- --testPathPattern="<component-name>.test.tsx"
```

Then run the full suite to confirm nothing regressed and coverage thresholds are met:

```bash
cd frontend && pnpm run test:coverage
```

Only mark the task complete when `pnpm run test:coverage` exits with code 0 and all thresholds pass.

---

## Rules

- **Always get approval before writing new tests** (Phase 4 gate)
- Never modify production component code solely to make a test pass — only fix genuine bugs
- If a component requires a Provider (e.g. Redux store, theme), wrap the render call in a test helper or inline wrapper — do not modify the component
- Coverage is collected only from `src/components/ui/**/*.{ts,tsx}` (see `jest.config.js`) — focus test additions there
