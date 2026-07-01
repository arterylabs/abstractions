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

export const metadata: Metadata = {
  title: "Abstractions — Keep your codebase knowledge fresh",
  description:
    "Receive weekly AI-generated newsletters that explain your codebase so your engineering team always knows how everything fits together. No more knowledge silos or onboarding bottlenecks.",
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
    icon: "lock_open",
    title: "Knowledge that never leaves your team",
    body: "Your codebase stays yours. Encrypted storage, zero training on your data, full control over who sees what.",
  },
  {
    icon: "sync",
    title: "Always in sync with reality",
    body: "No stale documentation. Abstractions watches your repos and updates automatically, so context stays fresh as your code evolves.",
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

const websiteSchema = {
  "@context": "https://schema.org",
  "@type": "WebSite",
  name: "Abstractions",
  url: siteUrl,
  description:
    "AI-generated newsletters that keep your engineering team in sync with your codebase.",
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
    "Tool that indexes your GitHub repositories and delivers weekly AI-generated newsletters to your engineering team.",
  url: siteUrl,
};

/* Eyebrow label — mono, uppercase, used above section headings */
function Eyebrow({ children }: { children: React.ReactNode }) {
  return (
    <span className="font-mono text-xs uppercase tracking-[0.2em] text-primary">
      {children}
    </span>
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
    <div className="gradient-border rounded-2xl p-1.5 shadow-2xl">
      <div className="overflow-hidden rounded-xl bg-card">
        {/* Masthead — matches the logo header of the real email */}
        <div className="border-b-[3px] border-primary">
          <div className="flex items-center justify-center gap-2 border-b border-border py-5">
            <Icon type="auto_awesome" className="!text-lg text-primary" />
            <span className="text-sm font-semibold tracking-tight text-foreground">
              Abstractions
            </span>
          </div>
        </div>

        <div className="flex flex-col gap-4 px-6 py-7">
          {/* Title (h1 in the email) */}
          <h3 className="text-lg font-bold leading-snug text-foreground">
            Idempotency keys moved to Redis
          </h3>

          {/* Article-style prose body */}
          <p className="text-sm leading-relaxed text-muted-foreground">
            Payment retries are now deduplicated at the gateway rather than in
            the database. The team moved idempotency keys into{" "}
            <code className="rounded border border-border bg-background/60 px-1 py-0.5 font-mono text-[0.72rem] text-foreground">
              Redis
            </code>{" "}
            to close the race condition surfaced during the March incident.
          </p>

          <h4 className="pt-1 text-[0.95rem] font-bold text-foreground">
            Why it matters
          </h4>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Duplicate charges were possible when a client retried before the
            first write committed. Centralising the key check removes that
            window without adding a second round-trip to Postgres.
          </p>

          {/* "About this newsletter" metadata block */}
          <div className="mt-3 flex flex-col gap-2 border-t border-border pt-5">
            <span className="text-[0.7rem] font-bold uppercase tracking-wide text-muted-foreground/70">
              About this newsletter
            </span>
            <span className="text-xs text-muted-foreground">
              <span className="font-medium text-foreground">Repository:</span>{" "}
              <span className="text-primary">payments-service</span>
            </span>
            <span className="flex flex-wrap items-center gap-1.5 text-xs text-muted-foreground">
              <span className="font-medium text-foreground">
                Files covered:
              </span>
              {["gateway/retry.go", "store/idempotency.go"].map((file) => (
                <code
                  key={file}
                  className="rounded border border-border bg-background/60 px-1.5 py-0.5 font-mono text-[0.68rem] text-foreground"
                >
                  {file}
                </code>
              ))}
            </span>
            <span className="text-[0.7rem] italic leading-relaxed text-muted-foreground/70">
              This content is AI-generated. Suggested improvements may not apply
              to your codebase — review them with that in mind.
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col bg-background text-foreground">
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
        <section className="relative overflow-hidden px-4 pb-24 pt-20 md:pt-28">
          {/* Texture + glow */}
          <div
            aria-hidden
            className="bg-grid mask-fade pointer-events-none absolute inset-0 -z-10 opacity-60"
          />
          <div
            aria-hidden
            className="hero-glow pointer-events-none absolute inset-x-0 -top-32 -z-10 h-[36rem]"
          />

          <div className="mx-auto flex max-w-[1200px] flex-col items-center gap-10">
            <div className="flex max-w-[760px] flex-col items-center gap-6 text-center">
              <span className="inline-flex items-center gap-2 rounded-full border border-border bg-card/60 px-3.5 py-1.5 text-xs font-medium text-muted-foreground backdrop-blur">
                <span className="relative flex size-1.5">
                  <span className="absolute inline-flex h-full w-full animate-glow-pulse rounded-full bg-primary" />
                  <span className="relative inline-flex size-1.5 rounded-full bg-primary" />
                </span>
                AI-generated newsletters for your codebase
              </span>

              <h1 className="text-balance text-4xl font-semibold leading-[1.08] tracking-tight md:text-5xl lg:text-6xl">
                <span className="text-gradient">
                  Retain codebase context as your systems evolve at agentic
                  speed.
                </span>
              </h1>

              <p className="max-w-[600px] text-pretty text-lg leading-relaxed text-muted-foreground">
                Abstractions indexes your repositories and ships weekly
                newsletters that explain architectural changes, system
                relationships, and emerging patterns — so your team stays
                aligned without knowledge silos or onboarding friction.
              </p>

              <div className="mt-1 flex flex-col items-center gap-3 sm:flex-row">
                <Button size="lg" href="/signup">
                  Start free
                  <Icon type="arrow_forward" className="!text-base" />
                </Button>
                <Button size="lg" variant="outline" href="/#how-it-works">
                  See how it works
                </Button>
              </div>

              <p className="font-mono text-xs text-muted-foreground">
                Bring your own Anthropic or OpenAI key · Your code is never
                stored
              </p>
            </div>

            <div className="w-full max-w-[640px]">
              <NewsletterPreview />
            </div>
          </div>
        </section>

        {/* ── Problem ──────────────────────────────────────────── */}
        <section id="problem" className="border-t border-border px-4 py-24">
          <div className="mx-auto flex max-w-[1200px] flex-col gap-14">
            <div className="flex flex-col items-center gap-4 text-center">
              <Eyebrow>The problem</Eyebrow>
              <h2 className="max-w-[640px] text-3xl font-semibold md:text-4xl">
                Context debt is silently killing your team&apos;s velocity
              </h2>
              <p className="max-w-[560px] leading-relaxed text-muted-foreground">
                The faster you ship, the harder it becomes for everyone on the
                team to understand how everything fits together.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-5 md:grid-cols-3">
              {painPoints.map((point) => (
                <div
                  key={point.title}
                  className="card-glow flex flex-col gap-4 rounded-xl border border-border bg-card p-6"
                >
                  <span className="flex size-11 items-center justify-center rounded-lg bg-accent text-muted-foreground">
                    <Icon type={point.icon} className="!text-xl" />
                  </span>
                  <h3 className="text-base font-semibold">{point.title}</h3>
                  <p className="text-sm leading-relaxed text-muted-foreground">
                    {point.body}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* ── Features ─────────────────────────────────────────── */}
        <section id="features" className="border-t border-border px-4 py-24">
          <div className="mx-auto flex max-w-[1200px] flex-col gap-14">
            <div className="flex flex-col items-center gap-4 text-center">
              <Eyebrow>Why Abstractions</Eyebrow>
              <h2 className="max-w-[640px] text-3xl font-semibold md:text-4xl">
                Solve context debt with the right tools
              </h2>
              <p className="max-w-[520px] leading-relaxed text-muted-foreground">
                Keep your entire team aligned with your codebase. Automatic,
                always fresh.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
              {features.map((feature) => (
                <div
                  key={feature.title}
                  className="card-glow flex flex-col gap-4 rounded-2xl border border-border bg-card p-8"
                >
                  <span className="flex size-12 items-center justify-center rounded-xl bg-primary/12 text-primary ring-1 ring-inset ring-primary/20">
                    <Icon type={feature.icon} className="!text-2xl" />
                  </span>
                  <h3 className="text-xl font-semibold">{feature.title}</h3>
                  <p className="leading-relaxed text-muted-foreground">
                    {feature.body}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* ── How it works ─────────────────────────────────────── */}
        <section
          id="how-it-works"
          className="relative overflow-hidden border-t border-border px-4 py-24"
        >
          <div
            aria-hidden
            className="bg-dots mask-fade pointer-events-none absolute inset-0 -z-10 opacity-50"
          />
          <div className="mx-auto flex max-w-[1200px] flex-col gap-14">
            <div className="flex flex-col items-center gap-4 text-center">
              <Eyebrow>How it works</Eyebrow>
              <h2 className="max-w-[640px] text-3xl font-semibold md:text-4xl">
                From commit history to inbox in four steps
              </h2>
              <p className="max-w-[560px] leading-relaxed text-muted-foreground">
                Abstractions connects to your GitHub repositories, deeply
                indexes your codebase, and delivers AI-generated newsletters on
                a schedule you control.
              </p>
            </div>

            <div className="relative grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {/* Connecting line on large screens */}
              <div
                aria-hidden
                className="absolute left-0 right-0 top-[2.4rem] -z-10 hidden h-px bg-gradient-to-r from-transparent via-border to-transparent lg:block"
              />
              {steps.map((step) => (
                <div
                  key={step.step}
                  className="flex flex-col gap-4 rounded-xl border border-border bg-card p-6"
                >
                  <div className="flex items-center justify-between">
                    <span className="flex size-11 items-center justify-center rounded-lg bg-primary/12 text-primary ring-1 ring-inset ring-primary/20">
                      <Icon type={step.icon} className="!text-xl" />
                    </span>
                    <span className="font-mono text-sm text-muted-foreground">
                      {step.step}
                    </span>
                  </div>
                  <h3 className="text-base font-semibold">{step.title}</h3>
                  <p className="text-sm leading-relaxed text-muted-foreground">
                    {step.body}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* ── FAQ ──────────────────────────────────────────────── */}
        <section id="faq" className="border-t border-border px-4 py-24">
          <div className="mx-auto flex max-w-[760px] flex-col gap-12">
            <div className="flex flex-col items-center gap-4 text-center">
              <Eyebrow>FAQ</Eyebrow>
              <h2 className="text-3xl font-semibold md:text-4xl">
                Frequently asked questions
              </h2>
              <p className="max-w-[500px] leading-relaxed text-muted-foreground">
                Everything you need to know about how Abstractions works.
              </p>
            </div>

            <Accordion
              type="single"
              collapsible
              className="flex flex-col gap-3"
            >
              {faqs.map((faq, i) => (
                <AccordionItem
                  key={faq.question}
                  value={`faq-${i}`}
                  className="rounded-xl border border-border bg-card px-5 last:border-b"
                >
                  <AccordionTrigger className="text-left text-base font-medium hover:no-underline">
                    {faq.question}
                  </AccordionTrigger>
                  <AccordionContent className="leading-relaxed text-muted-foreground">
                    {faq.answer}
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          </div>
        </section>

        {/* ── Final CTA ────────────────────────────────────────── */}
        <section className="px-4 pb-28 pt-4">
          <div className="mx-auto max-w-[1200px]">
            <div className="relative overflow-hidden rounded-3xl border border-border bg-card px-6 py-16 text-center">
              <div
                aria-hidden
                className="hero-glow pointer-events-none absolute inset-x-0 -bottom-24 top-0 -z-0 opacity-80"
              />
              <div className="relative z-10 mx-auto flex max-w-[560px] flex-col items-center gap-6">
                <h2 className="text-3xl font-semibold md:text-4xl">
                  Ready to keep your team in sync?
                </h2>
                <p className="leading-relaxed text-muted-foreground">
                  Connect your codebase and your first newsletter ships next
                  Monday.
                </p>
                <div className="flex flex-col items-center gap-3 sm:flex-row">
                  <Button size="lg" href="/signup">
                    Sign up free
                    <Icon type="arrow_forward" className="!text-base" />
                  </Button>
                  <Button size="lg" variant="outline" href="/docs">
                    Read the docs
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
