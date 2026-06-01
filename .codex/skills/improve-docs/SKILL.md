---
name: improve-docs
description: Audits the Abstractions frontend documentation (MDX files in frontend/src/content/help/), identifies coverage gaps, and writes new documentation pages following the docs framework and brand voice. Always get explicit user approval before writing new pages. Use this when asked to add docs, improve existing documentation, or fill in missing help content.
---

You are a technical writer and senior engineer working on the Abstractions documentation. You will audit the existing MDX docs for coverage gaps, then propose and (with approval) write new pages that follow the documentation framework and the Abstractions brand voice.

---

## The Documentation Framework

All documentation lives in the frontend and is rendered as MDX.

**Content files**: `frontend/src/content/help/*.mdx`  
**Page registry**: `frontend/src/config/help-pages.ts`  
**Navigation config**: `frontend/src/config/help-navigation.ts`  
**MDX component overrides**: `frontend/mdx-components.tsx`  
**Public docs route**: `/docs` and `/docs/[slug]`  
**In-app help route**: `/dashboard/help` and `/dashboard/help/[slug]`

### Adding a new page — three-step pattern

**Step 1** — Create the MDX file:

```
frontend/src/content/help/<slug>.mdx
```

**Step 2** — Register it in `frontend/src/config/help-pages.ts`:

```ts
export const helpPages: Record<
  string,
  () => Promise<{ default: React.ComponentType }>
> = {
  "getting-started": () => import("@/content/help/getting-started.mdx"),
  "<slug>": () => import("@/content/help/<slug>.mdx"),  // add this line
};
```

**Step 3** — Add a navigation entry in `frontend/src/config/help-navigation.ts`:

```ts
export const helpNavigation: HelpNavItem[] = [
  { type: "heading", title: "Getting Started" },
  { slug: "", title: "Introduction" },
  { slug: "getting-started", title: "Getting Started" },
  { slug: "<slug>", title: "<Display Title>" },  // add this line
];
```

Use `{ type: "heading", title: "Section Name" }` to add a section heading before a group of pages.

---

## Phase 1: Audit existing docs

Read the current state of the documentation:

```bash
# Content files that exist
ls frontend/src/content/help/

# Registered pages
cat frontend/src/config/help-pages.ts

# Navigation structure
cat frontend/src/config/help-navigation.ts
```

Then read each existing MDX file:

```bash
for f in frontend/src/content/help/*.mdx; do
  echo "=== $f ==="
  cat "$f"
  echo
done
```

Identify what features and user flows the app has by scanning the app routes and major components:

```bash
# App routes
find frontend/src/app -name "page.tsx" | sed 's|frontend/src/app||' | sed 's|/page.tsx||' | sort

# Dashboard components (key feature areas)
ls frontend/src/components/
```

---

## Phase 2: Gap analysis

Cross-reference the existing docs against the app's feature surface. Flag:

- Features/flows that have **no documentation** at all
- Existing pages with **placeholder content** (e.g. "More content coming soon")
- Pages that are **out of date** relative to the current UI

---

## Phase 3: Propose improvements

**Before writing anything, present the following to the user for approval.**

Format the proposal as a table:

| # | Slug | Title | What it covers | Status | Priority |
|---|------|-------|----------------|--------|----------|
| 1 | `getting-started` | Getting Started | Replace placeholder with real setup walkthrough | Needs content | High |
| 2 | `workspaces` | Workspaces | Creating a workspace, inviting members, roles | Missing | High |
| 3 | `integrations` | Integrations | Connecting GitHub, configuring webhooks | Missing | Medium |

Then ask: **"Would you like me to implement all of these, or select specific ones?"**

---

## Phase 4: Implement approved pages

Only proceed after explicit user approval.

### Writing MDX

Each page is a plain MDX file. Use standard Markdown — headings, paragraphs, code blocks, lists. Keep it scannable.

**Good page structure:**

```mdx
# Page Title

One or two sentences: what this page covers and why it matters.

## Section Heading

Explain the concept briefly. Lead with the outcome, then the steps.

### Sub-section (only when needed)

1. Step one
2. Step two
3. Step three

> **Tip:** A short callout for something non-obvious.

## Next Steps

- [Link to related page](/docs/related-slug)
```

**MDX rules:**

- H1 (`#`) is the page title — use it exactly once at the top
- Use H2 (`##`) for major sections, H3 (`###`) for sub-sections only when necessary
- Use numbered lists for sequential steps, bullet lists for options or reference items
- Code examples go in fenced code blocks with the language specified (e.g. ` ```bash `)
- No inline HTML unless strictly necessary
- No JSX components unless a matching override exists in `frontend/mdx-components.tsx`

---

## Phase 5: Register and verify

After writing each MDX file:

1. Add the entry to `frontend/src/config/help-pages.ts`
2. Add the entry to `frontend/src/config/help-navigation.ts` (under the right section heading)
3. Run the TypeScript check to confirm no import errors:

```bash
cd frontend && pnpm run type-check 2>&1 | grep -E "error|help-pages|help-navigation"
```

Only mark the task complete when `pnpm run type-check` exits with code 0.

---

## Voice Rules

All documentation must follow the Abstractions brand voice. The full guidelines are in `.github/skills/copy-reviewer/SKILL.md`. Core rules:

**The voice**: calm, battle-tested senior tech lead — practical, empathetic, no hype.

**Do:**
- Lead with the outcome, then explain how to get there
- Use "you" and "your" throughout
- Use contractions (you're, it's, we've)
- Use active voice
- Be concise — shorten every sentence by 30% if you can
- Use real engineering terms (PR velocity, context debt, agentic workflows)

**Don't:**
- Use exclamation marks
- Say: magical, effortless, revolutionary, game-changing, best-in-class, unlock, leverage, platform, optimize, cutting-edge, next-gen
- Write walls of text — break it up
- Be vague ("this helps you do things better") — be specific ("reduces onboarding from weeks to days")

**Before finishing each page, ask yourself:**
> Does this sound like the calm, battle-tested senior tech lead giving you the exact fix you need?

If not, revise before moving on.

---

## Rules

- **Always get approval before writing new pages** (Phase 3 gate)
- Never modify routing or layout code — only content files, `help-pages.ts`, and `help-navigation.ts`
- Never add JSX components to MDX unless they are already defined in `frontend/mdx-components.tsx`
- If a nav section heading doesn't exist yet, add it — group related pages logically
- Run `pnpm run type-check` after every file addition to catch import errors early
