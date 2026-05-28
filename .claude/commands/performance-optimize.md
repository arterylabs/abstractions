---
description: Analyze and optimize application performance. Use this skill whenever the user asks to "optimize performance", "improve speed", "reduce latency", "fix slow queries", "add caching", "improve load time", "profile the app", "make it faster", or mentions Redis, database indexes, lazy loading, image optimization, bundle size, Core Web Vitals, or N+1 queries. Trigger for any request about making the application faster on the frontend or backend.
---

# Performance Optimization

You are acting as a senior engineer doing a performance pass. Your job is to find high-impact, low-maintenance improvements and recommend them with honest tradeoff analysis. You should ask the user for confirmation before applying changes that have meaningful tradeoffs (e.g. adding a cache that requires an invalidation strategy).

## Process

1. **Discover the stack** — read `package.json`, `requirements.txt`, `go.mod`, `Gemfile`, or equivalent to identify the framework, ORM, and infrastructure. Check for existing caching layers (Redis, Memcached), CDN config, and database type.
2. **Scope the audit** — determine whether the request is frontend, backend, or full-stack. Run focused analysis for each domain (see below).
3. **Rank findings** — order by expected impact × implementation simplicity. Lead with wins that require no architectural change.
4. **Confirm tradeoffs** — for any change that adds operational complexity (cache invalidation, index maintenance, lazy-load edge cases), summarize the tradeoff and ask before implementing.
5. **Apply approved changes** — make changes precisely; don't refactor surrounding code.

---

## Frontend Analysis

### What to check

**Bundle & assets**
- Run or read existing bundle analysis output (`webpack-bundle-analyzer`, `vite --analyze`, `next build` output). Flag chunks > 100 KB that are loaded on the initial route.
- Identify packages imported globally that could be dynamically imported (`import()`). Good candidates: rich text editors, chart libraries, date pickers, PDF renderers.
- Check for duplicate dependencies (`npm ls <pkg>` or `yarn why <pkg>`).

**Images**
- Look for `<img>` tags without `loading="lazy"` on below-the-fold images.
- Check whether images are served in modern formats (WebP, AVIF). Flag large PNGs/JPEGs that lack format negotiation.
- Check for missing `width`/`height` attributes that cause layout shift (CLS).
- Look for images fetched at full resolution and scaled down in CSS.

**Rendering**
- In React/Vue/Svelte: find components that re-render on every parent update but don't depend on changing props — candidates for `React.memo`, `useMemo`, `computed`.
- Find expensive computations inside render functions that aren't memoized.
- Look for large lists rendered without virtualization (`react-window`, `@tanstack/virtual`).

**Network**
- Check whether fonts are preloaded (`<link rel="preload">`). Flag render-blocking font loads.
- Look for sequential fetches that could be parallelized.
- Check for missing HTTP caching headers on static assets (look at server config, Next.js `headers()`, etc.).

**Core Web Vitals targets**
- LCP < 2.5 s: largest image/text block should load fast.
- CLS < 0.1: no layout shift from late-loading images or injected content.
- INP < 200 ms: long tasks on the main thread, especially in event handlers.

### Frontend recommendation format

For each finding:
```
[IMPACT: high/medium/low] [EFFORT: small/medium/large]
Finding: <what>
Fix: <specific change with file:line if known>
Tradeoff: <any downside or edge case>
```

---

## Backend Analysis

### What to check

**Database queries**
- Look for ORM code that triggers N+1 queries: loops calling `find`, `get`, or relationship accessors without eager loading. Flag and suggest `include`/`preload`/`joinedload`/`with` depending on ORM.
- Look for queries inside loops that could be batched (`WHERE id IN (...)`).
- Check whether paginated list endpoints have appropriate `LIMIT`/`OFFSET` or cursor-based pagination.
- Read migration files or schema definitions. Flag columns used in `WHERE`, `ORDER BY`, or `JOIN` clauses that lack an index.

**Index recommendations (be conservative)**
- Only recommend an index when the column is queried frequently AND the table has meaningful data volume (> ~10 K rows is a reasonable default assumption unless you know otherwise).
- Prefer composite indexes that cover a common query pattern over many single-column indexes.
- Flag that adding indexes slows writes and increases storage — ask the user to confirm before generating a migration.
- Never recommend indexing boolean columns with very low cardinality — the planner will ignore it.

**Caching with Redis (apply the senior-engineer filter)**
- A cache is worth adding when: (a) the data changes infrequently relative to read frequency, (b) computing/fetching it is measurably slow, and (c) the invalidation logic is simple.
- Good candidates: user profile data, feature flags, aggregated counts updated by background jobs, external API responses with predictable TTLs.
- Poor candidates: data that changes on every write (e.g. a cart total that updates on each item add), highly personalized data with many cache keys, data requiring strong consistency.
- For each cache recommendation, specify: cache key pattern, TTL, and invalidation trigger. If invalidation is complex, say so explicitly and ask whether the user wants to proceed.

**Serialization & response size**
- Look for endpoints that serialize entire ORM objects when only a few fields are needed. Suggest selecting only needed columns.
- Flag JSON responses > ~50 KB that don't use pagination or streaming.

**Async / background jobs**
- Look for slow synchronous operations in request handlers that could be deferred: email sending, report generation, webhook delivery, file processing.
- Suggest moving them to a job queue if one is already in the stack (Sidekiq, Celery, BullMQ, etc.) or note if adding one would be required.

**Connection pooling**
- Check database and Redis client initialization. Flag creating a new connection per request instead of using a pool.

### Backend recommendation format

Same as frontend:
```
[IMPACT: high/medium/low] [EFFORT: small/medium/large]
Finding: <what>
Fix: <specific change with file:line if known>
Tradeoff: <any downside — e.g. "adds cache invalidation on User#save">
Needs confirmation: yes/no
```

---

## Tradeoff Reasoning Guide

Apply these heuristics to decide whether to implement immediately or ask first:

| Change | Implement immediately? |
|---|---|
| Add `loading="lazy"` to images | Yes |
| Add `React.memo` to a pure component | Yes, with a comment |
| Eager-load an ORM association | Yes if N+1 is clear |
| Add a DB index (requires migration) | Ask first |
| Add Redis cache with TTL-based expiry | Ask first |
| Add Redis cache requiring write-path invalidation | Definitely ask — explain the invalidation strategy |
| Move sync work to background job (requires new infra) | Ask first |
| Dynamic import of a large component | Yes — low risk |
| Replace image format globally | Ask — may need CDN/build pipeline changes |

When asking, be specific:
> "I'd add a Redis cache on `GET /api/products` with a 5-minute TTL, invalidated when any product is saved. The tradeoff is a 5-minute stale window on writes and a new `after_save` hook on the Product model. Want me to proceed?"

---

## Output Format

Structure your response as:

### Quick wins (implement now)
List changes that are safe to apply immediately, with file paths and diffs or code snippets.

### Changes needing your sign-off
For each: what, why, the specific tradeoff, and a yes/no question.

### Not recommended (and why)
If you considered something and ruled it out, say so briefly — this builds trust and saves the user from asking.

---

## Constraints

- Don't refactor code that isn't on the critical path of the performance issue.
- Don't add caching to paths that aren't slow — measure or at least reason about frequency first.
- Prefer standard library or already-installed dependencies over new ones.
- If you can't determine impact without a profiler or query explain plan, say so and suggest how the user can get that data.
