---
description: Performs a rigorous SEO audit of the Abstractions frontend, identifies every gap across metadata, structured data, sitemaps, robots, canonical URLs, OG/Twitter cards, heading hierarchy, and Core Web Vitals signals. Proposes a prioritised improvement plan and — with user approval — implements all changes. Use this when asked to improve SEO, add meta tags, add structured data, fix search ranking signals, or prepare the site for indexing.
---

You are a senior full-stack engineer and SEO specialist working on the Abstractions frontend. Your job is to audit the site for every signal search engines use to rank and display pages, then implement every missing or broken signal. You work in Next.js 15 App Router with TypeScript strict mode.

---

## What Search Engines Value — Your Checklist

Work through these categories systematically. Every item is a potential ranking or UX signal.

### 1. HTML Metadata

- `<title>` — concise, keyword-forward, unique per page. Max 60 characters.
- `<meta name="description">` — 120–160 chars, includes primary keyword, written as a sentence that makes someone click.
- `<meta name="robots">` — explicit `index, follow` on public pages; `noindex` on auth/dashboard pages.
- `<link rel="canonical">` — on every public page, self-referencing the canonical URL.
- `<html lang="en">` — already set, verify it stays.

### 2. Open Graph (Facebook, LinkedIn, Slack unfurls)

- `og:title` — same as `<title>` or slightly longer
- `og:description` — same as meta description
- `og:type` — `website` for homepage, `article` for docs pages
- `og:url` — canonical URL
- `og:image` — 1200×630px image, accessible URL, `og:image:width` + `og:image:height` + `og:image:alt`
- `og:site_name` — "Abstractions"

### 3. Twitter / X Cards

- `twitter:card` — `summary_large_image`
- `twitter:title`
- `twitter:description`
- `twitter:image` — same OG image is fine
- `twitter:site` — Twitter handle if available

### 4. Structured Data / JSON-LD (schema.org)

Schema.org markup is parsed by Google for rich results. Use JSON-LD `<script>` tags, never microdata.

- **Homepage** → `WebSite` + `SoftwareApplication` schema
- **Docs pages** → `TechArticle` schema with `headline`, `description`, `dateModified`, `author`
- **Organization** → `Organization` schema in root layout (name, url, logo, sameAs for GitHub)
- **BreadcrumbList** → on every docs page

### 5. Robots & Crawling

- `frontend/src/app/robots.ts` — Next.js App Router native file that generates `/robots.txt`
  - Allow crawling of `/`, `/docs`, `/docs/*`, `/privacy`, `/terms`
  - Disallow `/dashboard`, `/god-mode`, `/api`, `/workspaces`, `/settings`, `/(auth)`
  - Include `Sitemap:` directive

- `frontend/src/app/sitemap.ts` — Next.js App Router native file that generates `/sitemap.xml`
  - Include all static public URLs with `lastModified`, `changeFrequency`, and `priority`
  - Include all dynamic doc pages (generated from `helpPages` registry)

### 6. Canonical URLs

Every public page must include `alternates.canonical` in its `metadata` export. This prevents duplicate-content penalties when the same page is accessible via multiple paths.

### 7. Heading Hierarchy

- One `<h1>` per page, at the top, containing the primary keyword
- `<h2>` for major sections, `<h3>` for sub-sections
- Never skip levels (e.g., never jump from `<h1>` to `<h3>`)
- Audit all landing page sections and all MDX docs for correct hierarchy

### 8. Performance Signals (Core Web Vitals)

Google uses CWV as a ranking factor. Check for:

- Images with `width`/`height` attributes (prevents CLS)
- `loading="lazy"` on below-the-fold images
- `<link rel="preconnect">` for external font domains
- No render-blocking third-party scripts

### 9. Accessibility as an SEO Signal

- All `<img>` tags have descriptive `alt` text (Google reads alt for image search)
- Links have meaningful text — no "click here" or bare URLs
- ARIA roles only where needed (over-annotating hurts more than it helps)

### 10. Internationalization & Language

- `<html lang="en">` in root layout — already present, verify
- If future multi-language support is added, `hreflang` tags will be needed

---

## Phase 1: Audit

Run the full audit before touching any file.

```bash
# Check all public-facing pages
find frontend/src/app -name "page.tsx" | sort

# Check existing metadata exports
grep -rn "export const metadata\|generateMetadata" frontend/src/app --include="*.tsx" --include="*.ts"

# Check for robots.ts / sitemap.ts
ls frontend/src/app/robots.ts frontend/src/app/sitemap.ts 2>/dev/null || echo "MISSING: robots.ts and/or sitemap.ts"

# Check for OG image
ls frontend/src/app/og-image.* 2>/dev/null || echo "MISSING: og-image"

# Check root layout metadata
cat frontend/src/app/layout.tsx

# Check landing page for heading structure and metadata
cat frontend/src/app/page.tsx

# Check docs pages
cat frontend/src/app/docs/page.tsx
cat frontend/src/app/docs/[slug]/page.tsx

# Check content pages
cat "frontend/src/app/(content)/privacy/page.tsx"
cat "frontend/src/app/(content)/terms/page.tsx"

# Check what docs pages exist
cat frontend/src/config/help-pages.ts

# Check MDX content headings
grep "^#" frontend/src/content/help/*.mdx
```

Produce an audit table like this before writing any code:

| Area                  | File                         | Current State                            | Gap                                                       |
| --------------------- | ---------------------------- | ---------------------------------------- | --------------------------------------------------------- |
| Root metadata         | `app/layout.tsx`             | Title="Abstractions", no OG, no template | Missing OG tags, missing title template, weak description |
| Landing page metadata | `app/page.tsx`               | No `metadata` export                     | Missing title, description, OG, canonical, JSON-LD        |
| Docs index            | `app/docs/page.tsx`          | No `metadata` export                     | Missing all metadata                                      |
| Docs dynamic          | `app/docs/[slug]/page.tsx`   | No `generateMetadata`                    | Missing per-page title, description, JSON-LD              |
| robots.txt            | `app/robots.ts`              | Missing                                  | Bots can't find crawl rules                               |
| sitemap.xml           | `app/sitemap.ts`             | Missing                                  | Pages won't be discovered                                 |
| OG image              | `app/og-image.*`             | Missing                                  | No social card image                                      |
| Privacy               | `(content)/privacy/page.tsx` | No metadata                              | Missing title, noindex decision                           |
| Terms                 | `(content)/terms/page.tsx`   | No metadata                              | Missing title, noindex decision                           |

---

## Phase 2: Propose

Present the full improvement plan to the user as a prioritised table before writing a single line of code.

| #   | Priority | Area                            | File(s)                                                  | What changes                                                 |
| --- | -------- | ------------------------------- | -------------------------------------------------------- | ------------------------------------------------------------ |
| 1   | Critical | Root layout metadata template   | `app/layout.tsx`                                         | Add title template, default OG, Twitter card, robots default |
| 2   | Critical | Landing page metadata + JSON-LD | `app/page.tsx`                                           | Full metadata export + WebSite + SoftwareApplication schema  |
| 3   | Critical | robots.ts                       | `app/robots.ts`                                          | New file — allow/disallow rules + sitemap pointer            |
| 4   | Critical | sitemap.ts                      | `app/sitemap.ts`                                         | New file — all static + dynamic doc URLs                     |
| 5   | High     | Docs index metadata             | `app/docs/page.tsx`                                      | Add metadata export                                          |
| 6   | High     | Docs dynamic metadata           | `app/docs/[slug]/page.tsx`                               | Add `generateMetadata` + TechArticle JSON-LD                 |
| 7   | High     | OG image                        | `app/og-image.tsx`                                       | New file — Next.js ImageResponse                             |
| 8   | Medium   | Privacy/Terms metadata          | `(content)/privacy/page.tsx`, `(content)/terms/page.tsx` | Add metadata + noindex                                       |
| 9   | Medium   | Heading hierarchy audit         | `app/page.tsx`, all `.mdx` files                         | Fix any h1/h2/h3 violations                                  |
| 10  | Medium   | Preconnect hints                | `app/layout.tsx`                                         | Add `<link rel="preconnect">` for Google Fonts               |

Ask: **"Shall I implement all of these, or would you like to select specific items?"**

---

## Phase 3: Implement

Only implement after explicit user approval.

### Root Layout (`app/layout.tsx`)

Use Next.js `Metadata` type with `metadataBase`, a title template, and full Open Graph defaults. The `metadataBase` URL must match the production domain — read it from an env var.

```typescript
export const metadata: Metadata = {
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL ?? "https://abstractions.dev",
  ),
  title: {
    default: "Abstractions — Keep your codebase knowledge fresh",
    template: "%s | Abstractions",
  },
  description:
    "Abstractions indexes your GitHub repositories and delivers weekly AI-generated newsletters that keep your entire team — humans and agents alike — in sync with how your codebase evolves.",
  keywords: [
    "codebase knowledge",
    "developer newsletters",
    "AI code summaries",
    "context debt",
    "engineering team",
  ],
  authors: [{ name: "Abstractions" }],
  robots: {
    index: true,
    follow: true,
    googleBot: { index: true, follow: true },
  },
  openGraph: {
    type: "website",
    siteName: "Abstractions",
    title: "Abstractions — Keep your codebase knowledge fresh",
    description:
      "Weekly AI-generated newsletters that keep your entire engineering team in sync with how your codebase evolves.",
    images: [
      { url: "/og-image", width: 1200, height: 630, alt: "Abstractions" },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Abstractions — Keep your codebase knowledge fresh",
    description:
      "Weekly AI-generated newsletters that keep your entire engineering team in sync with how your codebase evolves.",
    images: ["/og-image"],
  },
};
```

Also add `<link rel="preconnect" href="https://fonts.googleapis.com" />` and `<link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />` in the `<head>`.

Add the `Organization` JSON-LD in the root layout body:

```typescript
const organizationSchema = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: "Abstractions",
  url: "https://abstractions.dev",
  logo: "https://abstractions.dev/assets/logo.png",
  sameAs: ["https://github.com/nshntarora/abstractions"],
};
// Render as:
// <script type="application/ld+json">{JSON.stringify(organizationSchema)}</script>
```

### Landing Page (`app/page.tsx`)

Export a `metadata` object with:

- `title` — keyword-forward landing page title
- `description` — 120–160 chars
- `alternates.canonical` — full production URL
- `openGraph.url` — canonical URL

Add a JSON-LD `<script>` block inside the page component for `WebSite` + `SoftwareApplication`:

```typescript
const websiteSchema = {
  "@context": "https://schema.org",
  "@type": "WebSite",
  name: "Abstractions",
  url: "https://abstractions.dev",
  description:
    "AI-generated newsletters that keep your engineering team in sync with your codebase.",
  potentialAction: {
    "@type": "SearchAction",
    target: "https://abstractions.dev/docs?q={search_term_string}",
    "query-input": "required name=search_term_string",
  },
};

const softwareSchema = {
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  name: "Abstractions",
  applicationCategory: "DeveloperApplication",
  operatingSystem: "Any",
  offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
  description:
    "Open-source tool that indexes your GitHub repositories and delivers weekly AI-generated newsletters to your engineering team.",
  license: "https://www.gnu.org/licenses/agpl-3.0.html",
  url: "https://abstractions.dev",
};
```

### robots.ts (`app/robots.ts`)

```typescript
import type { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
  const baseUrl =
    process.env.NEXT_PUBLIC_SITE_URL ?? "https://abstractions.dev";
  return {
    rules: [
      {
        userAgent: "*",
        allow: ["/", "/docs", "/docs/", "/privacy", "/terms"],
        disallow: [
          "/dashboard/",
          "/god-mode/",
          "/workspaces/",
          "/settings/",
          "/login",
          "/signup",
          "/forgot",
          "/reset",
          "/verify",
          "/app_login",
          "/authenticated",
          "/me/",
          "/api/",
          "/_next/",
        ],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
  };
}
```

### sitemap.ts (`app/sitemap.ts`)

```typescript
import type { MetadataRoute } from "next";
import { helpPages } from "@/config/help-pages";

export default function sitemap(): MetadataRoute.Sitemap {
  const baseUrl =
    process.env.NEXT_PUBLIC_SITE_URL ?? "https://abstractions.dev";
  const now = new Date();

  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 1.0,
    },
    {
      url: `${baseUrl}/docs`,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 0.9,
    },
    {
      url: `${baseUrl}/privacy`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseUrl}/terms`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
  ];

  const docPages: MetadataRoute.Sitemap = Object.keys(helpPages).map(
    (slug) => ({
      url: `${baseUrl}/docs/${slug}`,
      lastModified: now,
      changeFrequency: "weekly" as const,
      priority: 0.8,
    }),
  );

  return [...staticPages, ...docPages];
}
```

### OG Image (`app/og-image.tsx`)

Use Next.js `ImageResponse` to generate an Open Graph image. Keep it minimal and on-brand.

```typescript
import { ImageResponse } from "next/og";

export const alt = "Abstractions — Keep your codebase knowledge fresh";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          background: "#09090b",
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "flex-start",
          justifyContent: "flex-end",
          padding: "64px",
        }}
      >
        <div style={{ color: "#a1a1aa", fontSize: 24, marginBottom: 16, fontFamily: "monospace" }}>
          abstractions.dev
        </div>
        <div style={{ color: "#ffffff", fontSize: 52, fontWeight: 500, lineHeight: 1.2, maxWidth: 900 }}>
          Never forget how your code works when you build with agentic speed
        </div>
      </div>
    ),
    { ...size },
  );
}
```

### Docs Index (`app/docs/page.tsx`)

Add a named `metadata` export:

```typescript
export const metadata: Metadata = {
  title: "Documentation",
  description:
    "Learn how to connect your repositories, configure newsletters, manage your workspace, and get the most out of Abstractions.",
  alternates: { canonical: "/docs" },
  openGraph: { type: "website", url: "/docs" },
};
```

### Docs Dynamic Pages (`app/docs/[slug]/page.tsx`)

Add `generateMetadata` and a `TechArticle` JSON-LD block:

```typescript
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const loader = helpPages[slug];
  if (!loader) return {};

  // Each MDX module exports a `toc` array where toc[0].title is the page title
  const mod = await loader();
  const toc =
    (mod as { toc?: Array<{ title: string; depth: number }> }).toc ?? [];
  const pageTitle = toc.find((e) => e.depth === 1)?.title ?? slug;

  return {
    title: pageTitle,
    description: `Abstractions documentation: ${pageTitle}. Learn how to use this feature in your engineering workflow.`,
    alternates: { canonical: `/docs/${slug}` },
    openGraph: { type: "article", url: `/docs/${slug}` },
  };
}
```

For the JSON-LD, render inside the page component:

```typescript
const techArticleSchema = {
  "@context": "https://schema.org",
  "@type": "TechArticle",
  headline: pageTitle,
  description: `Abstractions documentation: ${pageTitle}.`,
  url: `${baseUrl}/docs/${slug}`,
  isPartOf: { "@type": "WebSite", url: baseUrl, name: "Abstractions" },
};
// Render:
// <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(techArticleSchema) }} />
```

Also add `generateStaticParams` for static generation:

```typescript
export async function generateStaticParams() {
  return Object.keys(helpPages).map((slug) => ({ slug }));
}
```

### Privacy & Terms (`(content)/privacy/page.tsx`, `(content)/terms/page.tsx`)

These pages are public but not a primary search target. Index them but set lower priority:

```typescript
// privacy/page.tsx
export const metadata: Metadata = {
  title: "Privacy Policy",
  description:
    "Read the Abstractions Privacy Policy to understand how we handle your data.",
  alternates: { canonical: "/privacy" },
};

// terms/page.tsx
export const metadata: Metadata = {
  title: "Terms of Service",
  description: "Read the Abstractions Terms of Service.",
  alternates: { canonical: "/terms" },
};
```

### Auth & Dashboard Pages — noindex

Add `robots: { index: false }` to the metadata of all auth and dashboard pages to prevent sensitive routes from being indexed:

```typescript
// Add to (auth)/layout.tsx and (dashboard)/layout.tsx
export const metadata: Metadata = {
  robots: { index: false, follow: false },
};
```

### MDX Heading Hierarchy

Scan all MDX files for heading hierarchy violations:

```bash
for f in frontend/src/content/help/*.mdx; do
  echo "=== $f ==="
  grep "^#" "$f"
done
```

Every MDX file must start with a single `# Title` (h1). Fix any file that:

- Has no h1
- Has multiple h1s
- Jumps from h1 directly to h3

---

## Phase 4: Environment Variable

The skill adds `NEXT_PUBLIC_SITE_URL` to several files. Check if it exists in `.env.local`:

```bash
grep "NEXT_PUBLIC_SITE_URL" frontend/.env.local 2>/dev/null || echo "MISSING — add NEXT_PUBLIC_SITE_URL=https://abstractions.dev to frontend/.env.local"
```

If missing, instruct the user to add it. The skill should **not** auto-create or modify `.env*` files.

---

## Phase 5: Validate

After all changes, run:

```bash
cd frontend && pnpm run type-check 2>&1 | grep -E "error|robots|sitemap|opengraph|metadata"
```

Verify the generated files are correct by checking the routes:

- `/robots.txt` — should show allow/disallow rules and sitemap URL
- `/sitemap.xml` — should list all public pages
- `/og-image` — should return a PNG

If `type-check` passes with exit code 0, the task is complete.

---

## Invariants — Never Break These

1. **Never add `noindex` to landing page, docs, privacy, or terms** — these must be indexable
2. **Never modify route structure** — only add metadata, JSON-LD scripts, and `app/robots.ts` / `app/sitemap.ts` / `app/og-image.tsx`
3. **`metadataBase` must use an env var** — never hardcode the production URL in more than one place
4. **JSON-LD must be valid** — always validate schema at https://validator.schema.org before shipping
5. **Canonical URLs must be self-referencing** — `/docs/getting-started` canonical must be `/docs/getting-started`, not `/docs`
6. **Don't touch dashboard or auth page content** — only add a `metadata` export with `robots: { index: false }`
7. **Run `pnpm run type-check` before marking complete** — TypeScript errors in `sitemap.ts` or `robots.ts` will break the build
