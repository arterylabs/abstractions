import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import LandingPageNavbar from "@/components/landing/Navbar";
import Footer from "@/components/landing/Footer";
import { ogImages, siteUrl } from "./shared-metadata";

const GITHUB_URL = "https://github.com/arterylabs/abstractions";
const LICENSE_URL =
  "https://github.com/arterylabs/abstractions/blob/master/LICENSE";

export const metadata: Metadata = {
  title: "Abstractions — Receive weekly newsletters explaining your codebase",
  description:
    "Open-source tool that indexes your GitHub repositories and delivers weekly AI-generated newsletters — without you having to ask — so your engineering team always knows how everything fits together. AGPL-3.0 — use the hosted product or self-host.",
  alternates: {
    canonical: siteUrl,
  },
  openGraph: {
    type: "website",
    url: siteUrl,
    images: ogImages,
  },
};

const painPoints = [
  {
    icon: "person_search",
    title: "New hires can't keep up",
    body: "PRs merge faster than onboarding docs can keep pace. New engineers spend weeks just building a mental model.",
  },
  {
    icon: "memory",
    title: "Senior engineers start forgetting",
    body: "How did we handle auth again? Which service owns rate-limiting? Context gets lost as codebases evolve.",
  },
  {
    icon: "history_edu",
    title: "Decisions get made without history",
    body: "Why did we pick this pattern? What trade-offs did we consider? That reasoning lives in someone's head — or nowhere at all.",
  },
];

const features = [
  {
    icon: "lightbulb",
    title: "Understand the why, not just the what",
    body: "Newsletters explain the decisions, trade-offs, and reasoning behind your code — not just what it does, but why it was built that way.",
  },
  {
    icon: "speed",
    title: "Onboard engineers in days, not weeks",
    body: "New hires get instant access to an AI-generated guide to your codebase. Understand patterns, decisions, and dependencies on day one.",
  },
  {
    icon: "sync",
    title: "Doesn't wait for you to ask",
    body: "You could ask an AI assistant about your codebase — but that only surfaces what you think to ask about. Abstractions watches your repos, figures out a topic on its own, and tells you before you knew to ask.",
  },
  {
    icon: "code",
    title: "Open source, run it your way",
    body: "Licensed under AGPL-3.0. Inspect every line, self-host on your own infrastructure, or use the hosted product — same codebase either way. Your data stays yours.",
  },
];

const steps = [
  {
    icon: "cable",
    step: "01",
    title: "Connect GitHub",
    body: "Authorize your repositories — public or private — with a scoped OAuth connection.",
  },
  {
    icon: "manage_search",
    step: "02",
    title: "Index the codebase",
    body: "We analyze commits and PRs, chunk the code, and map concepts into a topic graph.",
  },
  {
    icon: "auto_awesome",
    step: "03",
    title: "Generate the newsletter",
    body: "Your model of choice drafts a readable changelog from the indexed context.",
  },
  {
    icon: "mark_email_read",
    step: "04",
    title: "Deliver to the team",
    body: "Ships to your engineers on the schedule you control — or on demand, anytime.",
  },
];

const faqs = [
  {
    question:
      "Why can't I just ask my coding agent to generate these newsletters?",
    answer:
      "You can — but that puts the burden on you to know what to ask. Abstractions works in the background, surfacing knowledge about parts of your codebase that you or your teammates may have never touched. Instead of waiting for someone to think to ask, we deliver those insights automatically, covering areas of the codebase that often go undiscovered.",
  },
  {
    question: "What do you do with my code?",
    answer:
      "We process your code solely to generate summaries and embeddings used for newsletter creation. Once a summary is produced, it is not retained — we never persist your raw source code in our database. This approach gives us what we need to generate meaningful, accurate newsletters while ensuring your code never sits in storage.",
  },
  {
    question: "Is Abstractions open source?",
    answer:
      "Yes. The entire product — Go API, Next.js frontend, and workers — is open source under the GNU Affero General Public License v3.0 (AGPL-3.0). The hosted service at useabstractions.com runs the same codebase you can read, fork, and self-host on GitHub.",
  },
  {
    question: "Can I self-host Abstractions?",
    answer:
      "Yes. Clone the repository, configure Postgres, Redis, and your environment, and run the API, workers, and frontend yourself. Setup details live in the project README on GitHub. Note that trademarks and hosted branding for useabstractions.com are not licensed under the AGPL — if you run your own instance, plan to rebrand.",
  },
  {
    question: "Who will receive the newsletters?",
    answer:
      "You decide. Recipients are fully configurable at the repository level, so you can send to your entire engineering team, a specific sub-team, or any distribution list you choose. Most teams send to something like engineering@yourcompany.com or a team-specific alias.",
  },
  {
    question: "How do you decide which topics to write about?",
    answer:
      "Topics are automatically derived from themes and patterns we discover while indexing your repository. As your codebase evolves, so do the topics we surface. You can also manually suggest topics for upcoming newsletters directly within the repository settings, giving you full control when there is something specific you want the team to understand.",
  },
];

/* Sample of the raw material a newsletter is assembled from — scrolled as a
   wire-service ticker between the hero and the body of the page. */
const tickerItems = [
  "feat(gateway): move idempotency keys to Redis",
  "refactor(auth): collapse session + token middleware",
  "fix(billing): retry webhook on 5xx with jitter",
  "chore(db): add pgvector index on chunk embeddings",
  "feat(workers): fan out indexing per repository",
  "perf(search): cache topic graph between digests",
  "fix(email): suppress bounced recipients for 30d",
  "feat(api): scope OAuth to selected repositories",
];

const websiteSchema = {
  "@context": "https://schema.org",
  "@type": "WebSite",
  name: "Abstractions",
  url: siteUrl,
  description:
    "Open-source AI-generated newsletters that keep your engineering team in sync with your codebase.",
  potentialAction: {
    "@type": "SearchAction",
    target: `${siteUrl}/docs?q={search_term_string}`,
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
    "Open-source tool that indexes your GitHub repositories and delivers weekly AI-generated newsletters to your engineering team. AGPL-3.0 licensed.",
  url: siteUrl,
  license: LICENSE_URL,
  codeRepository: GITHUB_URL,
};

/* Section marker set as a folio: § number, running head, and a rule that
   carries the eye across the page. */
function SectionMark({ index, label }: { index: string; label: string }) {
  return (
    <div className="flex items-center gap-4">
      <span className="press-label text-primary">§ {index}</span>
      <span className="press-label text-muted-foreground">{label}</span>
      <span aria-hidden className="h-px flex-1 bg-border" />
    </div>
  );
}

/* Stylized preview of an actual newsletter email shown in the hero.
   Mirrors the real email structure the app sends (see
   api/clients/email/generic_email_template_html.tmpl and
   api/services/newsletter/assembler.go): a branded masthead, a title,
   long-form prose that explains the "why" behind changes, and an
   "About this newsletter" metadata block with the AI-generated disclaimer. */
function NewsletterPreview() {
  return (
    <div className="press-sheet relative rounded-sm border border-border bg-card">
      {/* Masthead — matches the logo header of the real email */}
      <div className="flex items-center justify-between border-b border-border px-6 py-3.5">
        <span className="press-label text-muted-foreground">Abstractions</span>
        <span className="press-label text-muted-foreground">Issue 14</span>
      </div>

      <div className="flex flex-col gap-4 px-6 py-7 sm:px-8">
        <span className="press-label text-primary">payments-service</span>

        {/* Title (h1 in the email) */}
        <h3 className="press-display text-[2rem] leading-[1.02]">
          Idempotency keys moved to Redis
        </h3>

        {/* Article-style prose body */}
        <p className="text-[0.9rem] leading-[1.75] text-muted-foreground">
          Payment retries are now deduplicated at the gateway rather than in the
          database. The team moved idempotency keys into{" "}
          <code className="border-b border-primary/50 font-mono text-[0.8rem] text-foreground">
            Redis
          </code>{" "}
          to close the race condition surfaced during the March incident.
        </p>

        <h4 className="press-display pt-1 text-xl">Why it matters</h4>
        <p className="text-[0.9rem] leading-[1.75] text-muted-foreground">
          Duplicate charges were possible when a client retried before the first
          write committed. Centralising the key check removes that window
          without adding a second round-trip to Postgres.
        </p>

        {/* "About this newsletter" metadata block */}
        <div className="mt-3 flex flex-col gap-2.5 border-t border-border pt-5">
          <span className="press-label text-muted-foreground/70">
            About this newsletter
          </span>
          <span className="flex flex-wrap items-center gap-1.5 font-mono text-[0.7rem] text-muted-foreground">
            Files covered:
            {["gateway/retry.go", "store/idempotency.go"].map((file) => (
              <code
                key={file}
                className="rounded-sm bg-accent px-1.5 py-0.5 text-[0.68rem] text-foreground"
              >
                {file}
              </code>
            ))}
          </span>
          <span className="text-[0.72rem] italic leading-relaxed text-muted-foreground/70">
            This content is AI-generated. Suggested improvements may not apply
            to your codebase — review them with that in mind.
          </span>
        </div>
      </div>
    </div>
  );
}

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(websiteSchema) }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(softwareSchema) }}
      />
      <LandingPageNavbar />

      <main>
        {/* ── Hero ─────────────────────────────────────────────── */}
        <section className="press-grain relative overflow-hidden px-6 pb-20 pt-10">
          <div
            aria-hidden
            className="press-wash mask-fade pointer-events-none absolute inset-x-0 -top-40 -z-10 h-[42rem]"
          />

          <div className="mx-auto max-w-[1180px]">
            {/* Masthead rule — the paper's dateline */}
            <div
              className="reveal flex items-center justify-between gap-4 border-y border-border py-2.5"
              style={{ "--d": "0ms" } as React.CSSProperties}
            >
              <span className="press-label text-muted-foreground">
                Vol. 01 — Weekly dispatch
              </span>
              <span className="press-label hidden text-muted-foreground sm:inline">
                Written by machines, for the people who ship
              </span>
              <span className="press-label flex items-center gap-2 text-primary">
                <span className="relative flex size-1">
                  <span
                    aria-hidden
                    className="animate-ink-pulse absolute inline-flex size-1 rounded-full bg-primary"
                  />
                  <span className="relative inline-flex size-1 rounded-full bg-primary" />
                </span>
                Open source
              </span>
            </div>

            <div className="grid grid-cols-1 items-start gap-14 pt-14 lg:grid-cols-[1.05fr_0.95fr] lg:gap-16 lg:pt-16">
              <div className="flex flex-col gap-7">
                <h1
                  className="press-display reveal text-balance text-[clamp(2.9rem,7.2vw,4.9rem)]"
                  style={{ "--d": "80ms" } as React.CSSProperties}
                >
                  Weekly newsletters that{" "}
                  <em className="italic text-primary">document</em> your
                  codebase
                </h1>

                <p
                  className="reveal max-w-[46ch] text-pretty text-lg leading-[1.7] text-muted-foreground"
                  style={{ "--d": "180ms" } as React.CSSProperties}
                >
                  You could ask an AI about your codebase — Abstractions just
                  doesn&apos;t wait for you to. Weekly digests of architectural
                  changes, system relationships, and patterns, delivered
                  automatically.
                </p>

                <div
                  className="reveal flex flex-col items-stretch gap-3 sm:flex-row sm:items-center"
                  style={{ "--d": "260ms" } as React.CSSProperties}
                >
                  <Button size="lg" href="/signup">
                    Start free
                    <Icon type="arrow_forward" className="!text-base" />
                  </Button>
                  <Button size="lg" variant="outline" asChild>
                    <a
                      href={GITHUB_URL}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      View on GitHub
                      <Icon type="open_in_new" className="!text-base" />
                    </a>
                  </Button>
                </div>

                {/* Colophon — the terms of the deal, set like small print */}
                <dl
                  className="reveal grid grid-cols-1 gap-6 border-t border-border pt-6 sm:grid-cols-3"
                  style={{ "--d": "340ms" } as React.CSSProperties}
                >
                  {[
                    ["License", "AGPL-3.0"],
                    ["Models", "Bring your own key"],
                    ["Your source", "Never stored"],
                  ].map(([term, value]) => (
                    <div key={term} className="flex flex-col gap-1.5">
                      <dt className="press-label text-muted-foreground/70">
                        {term}
                      </dt>
                      <dd className="font-mono text-sm text-foreground">
                        {value}
                      </dd>
                    </div>
                  ))}
                </dl>
              </div>

              <div
                className="reveal w-full"
                style={{ "--d": "420ms" } as React.CSSProperties}
              >
                <NewsletterPreview />
              </div>
            </div>
          </div>
        </section>

        {/* ── Wire ticker ──────────────────────────────────────── */}
        <div className="ticker-mask overflow-hidden border-y border-border bg-accent-dark py-3">
          <div className="ticker-track" aria-hidden>
            {[...tickerItems, ...tickerItems].map((item, i) => (
              <span
                key={i}
                className="flex shrink-0 items-center gap-6 whitespace-nowrap px-6 font-mono text-xs text-muted-foreground"
              >
                <span className="text-primary">◆</span>
                {item}
              </span>
            ))}
          </div>
        </div>

        {/* ── Problem ──────────────────────────────────────────── */}
        <section id="problem" className="px-6 py-24">
          <div className="mx-auto flex max-w-[1180px] flex-col gap-12">
            <SectionMark index="01" label="The problem" />

            <div className="grid grid-cols-1 gap-10 lg:grid-cols-[0.9fr_1.1fr] lg:gap-16">
              <h2 className="press-display max-w-[14ch] text-[clamp(2.1rem,4.2vw,3.1rem)]">
                Context debt is silently killing your team&apos;s velocity
              </h2>
              <p className="max-w-[52ch] self-end text-lg leading-[1.75] text-muted-foreground">
                The faster you ship, the harder it becomes for everyone on the
                team to understand how everything fits together.
              </p>
            </div>

            {/* Newspaper columns, separated by rules rather than card chrome */}
            <div className="grid grid-cols-1 border-t border-border md:grid-cols-3">
              {painPoints.map((point) => (
                <div
                  key={point.title}
                  className="flex flex-col gap-4 border-b border-border px-0 py-8 md:border-b-0 md:border-l md:px-8 md:first:border-l-0 md:first:pl-0 md:last:pr-0"
                >
                  <Icon type={point.icon} className="!text-xl text-primary" />
                  <h3 className="press-display text-2xl">{point.title}</h3>
                  <p className="drop-cap text-[0.95rem] leading-[1.75] text-muted-foreground">
                    {point.body}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* ── Features ─────────────────────────────────────────── */}
        <section id="features" className="px-6 py-24">
          <div className="mx-auto flex max-w-[1180px] flex-col gap-12">
            <SectionMark index="02" label="Why Abstractions" />

            <div className="grid grid-cols-1 gap-10 lg:grid-cols-[0.9fr_1.1fr] lg:gap-16">
              <h2 className="press-display max-w-[14ch] text-[clamp(2.1rem,4.2vw,3.1rem)]">
                Solve context debt with the right tools
              </h2>
              <p className="max-w-[52ch] self-end text-lg leading-[1.75] text-muted-foreground">
                Keep your entire team aligned with your codebase. Automatic,
                always fresh.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              {features.map((feature, i) => (
                <article
                  key={feature.title}
                  className="press-card flex flex-col gap-4 overflow-hidden border border-border bg-card p-8"
                >
                  <div className="flex items-baseline justify-between gap-4">
                    <span className="press-display text-3xl text-primary">
                      {String(i + 1).padStart(2, "0")}
                    </span>
                    <Icon
                      type={feature.icon}
                      className="!text-lg text-muted-foreground"
                    />
                  </div>
                  <h3 className="press-display max-w-[18ch] text-[1.7rem] leading-[1.1]">
                    {feature.title}
                  </h3>
                  <p className="text-[0.95rem] leading-[1.75] text-muted-foreground">
                    {feature.body}
                  </p>
                </article>
              ))}
            </div>
          </div>
        </section>

        {/* ── How it works ─────────────────────────────────────── */}
        <section
          id="how-it-works"
          className="relative overflow-hidden border-y border-border px-6 py-24"
        >
          <div
            aria-hidden
            className="bg-hatch mask-fade pointer-events-none absolute inset-0 -z-10 opacity-70"
          />
          <div className="mx-auto flex max-w-[1180px] flex-col gap-12">
            <SectionMark index="03" label="How it works" />

            <div className="grid grid-cols-1 gap-10 lg:grid-cols-[0.9fr_1.1fr] lg:gap-16">
              <h2 className="press-display max-w-[14ch] text-[clamp(2.1rem,4.2vw,3.1rem)]">
                From commit history to inbox in four steps
              </h2>
              <p className="max-w-[52ch] self-end text-lg leading-[1.75] text-muted-foreground">
                Abstractions connects to your GitHub repositories, deeply
                indexes your codebase, and delivers AI-generated newsletters on
                a schedule you control.
              </p>
            </div>

            {/* The press run — steps laid out along a single rule */}
            <ol className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
              {steps.map((step) => (
                <li
                  key={step.step}
                  className="group flex flex-col gap-4 border-t-2 border-border pt-6 transition-colors duration-300 hover:border-primary lg:pr-8"
                >
                  <div className="flex items-center gap-3">
                    <span className="font-mono text-sm text-primary">
                      {step.step}
                    </span>
                    <Icon
                      type={step.icon}
                      className="!text-lg text-muted-foreground transition-colors duration-300 group-hover:text-foreground"
                    />
                  </div>
                  <h3 className="press-display text-2xl">{step.title}</h3>
                  <p className="pb-8 text-[0.95rem] leading-[1.75] text-muted-foreground">
                    {step.body}
                  </p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* ── FAQ ──────────────────────────────────────────────── */}
        <section id="faq" className="px-6 py-24">
          <div className="mx-auto flex max-w-[1180px] flex-col gap-12">
            <SectionMark index="04" label="FAQ" />

            <div className="grid grid-cols-1 gap-10 lg:grid-cols-[0.9fr_1.1fr] lg:gap-16">
              <div className="flex flex-col gap-5">
                <h2 className="press-display max-w-[12ch] text-[clamp(2.1rem,4.2vw,3.1rem)]">
                  Frequently asked questions
                </h2>
                <p className="max-w-[38ch] leading-[1.75] text-muted-foreground">
                  Everything you need to know about how Abstractions works.
                  Anything else — the answer is probably in the{" "}
                  <Link href="/docs" className="underline underline-offset-4">
                    docs
                  </Link>
                  .
                </p>
              </div>

              <Accordion type="single" collapsible className="flex flex-col">
                {faqs.map((faq, i) => (
                  <AccordionItem
                    key={faq.question}
                    value={`faq-${i}`}
                    className="border-b border-border last:border-b"
                  >
                    <AccordionTrigger className="py-5 text-left font-[family-name:var(--font-heading)] text-xl font-normal transition-colors hover:text-primary hover:no-underline [&[data-state=open]>.icon]:rotate-180 [&[data-state=open]]:text-primary">
                      {faq.question}
                    </AccordionTrigger>
                    <AccordionContent className="max-w-[58ch] pb-6 text-[0.95rem] leading-[1.8] text-muted-foreground">
                      {faq.answer}
                    </AccordionContent>
                  </AccordionItem>
                ))}
              </Accordion>
            </div>
          </div>
        </section>

        {/* ── Final CTA ────────────────────────────────────────── */}
        <section className="px-6 pb-24">
          <div className="mx-auto max-w-[1180px]">
            <div className="press-grain relative overflow-hidden border border-border bg-card px-6 py-20 text-center sm:px-12">
              <div
                aria-hidden
                className="press-wash pointer-events-none absolute inset-0 -z-0 opacity-90"
              />
              <div className="relative z-10 mx-auto flex max-w-[38ch] flex-col items-center gap-6">
                <span className="press-label text-primary">
                  Next issue ships Monday
                </span>
                <h2 className="press-display text-[clamp(2.2rem,5vw,3.4rem)]">
                  Ready to keep your team in sync?
                </h2>
                <p className="leading-[1.75] text-muted-foreground">
                  Sign up for the hosted product, or clone the open-source repo
                  and run it yourself. Your first newsletter can ship next
                  Monday.
                </p>
                <div className="flex flex-col items-stretch gap-3 sm:flex-row sm:items-center">
                  <Button size="lg" href="/signup">
                    Sign up free
                    <Icon type="arrow_forward" className="!text-base" />
                  </Button>
                  <Button size="lg" variant="outline" asChild>
                    <a
                      href={GITHUB_URL}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      View on GitHub
                      <Icon type="open_in_new" className="!text-base" />
                    </a>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
