import Link from "next/link";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";
import LandingPageNavbar from "@/components/landing/Navbar";
import Footer from "@/components/landing/Footer";

const headingClass = "font-[family-name:var(--font-heading)]";

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

export default function LandingPage() {
  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground">
      <LandingPageNavbar />

      {/* Hero */}
      <section className="py-24 px-4">
        <div className="max-w-[900px] mx-auto flex flex-col items-center text-center gap-6">
          {/* Open source badge — hidden until open source launch */}
          {/* <div className="inline-flex items-center gap-2 border border-primary/40 text-primary text-sm px-3 py-1 rounded-full">
            <span className="size-1.5 rounded-full bg-primary inline-block" />
            Open Source · AGPL License
          </div> */}

          <h1
            className={`${headingClass} text-2xl md:text-3xl lg:text-4xl font-medium leading-tight`}
          >
            Retain your codebase context when it evolves with agentic speed
          </h1>

          <p className="text-lg text-muted-foreground max-w-[600px] leading-relaxed">
            Receive weekly newsletters that explain concepts in your codebase so
            your team always knows how everything fits together. No more
            knowledge silos, tribal knowledge, or onboarding bottlenecks.
          </p>

          <div className="flex flex-col sm:flex-row items-center gap-3 mt-2">
            <Link href="/signup">
              <Button size="lg">Sign up</Button>
            </Link>
            {/* View on GitHub — hidden until open source launch */}
            {/* <Link
              href="https://github.com/nshntarora/abstractions"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Button size="lg" variant="outline">
                View on GitHub
              </Button>
            </Link> */}
          </div>

          {/* AGPL license text — hidden until open source launch */}
          {/* <p className="text-sm text-muted-foreground">
            AGPL licensed. Audit the code yourself.
          </p> */}

          {/* Screenshot placeholder */}
          {/* <div className="w-full mt-8 bg-card border border-border rounded-xl aspect-video flex items-center justify-center">
            <p className="text-muted-foreground text-sm">
              App screenshot coming soon
            </p>
          </div> */}
        </div>
      </section>

      {/* Problem Statement */}
      <section id="problem" className="py-24 px-4 bg-card">
        <div className="max-w-[900px] mx-auto flex flex-col gap-12">
          <div className="flex flex-col items-center text-center gap-4">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              Context debt is silently killing your team&apos;s velocity
            </h2>
            <p className="text-muted-foreground max-w-[560px] leading-relaxed">
              The faster you ship, the harder it becomes for everyone on the
              team to understand how everything fits together.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {painPoints.map((point) => (
              <div
                key={point.title}
                className="bg-background border border-border rounded-xl p-6 flex flex-col gap-3"
              >
                <Icon type={point.icon} className="text-[2rem] text-primary" />
                <h3 className="font-medium text-base">{point.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {point.body}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="py-24 px-4">
        <div className="max-w-[900px] mx-auto flex flex-col gap-16">
          <div className="flex flex-col items-center text-center gap-4">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              Solve context debt with the right tools
            </h2>
            <p className="text-muted-foreground max-w-[500px] leading-relaxed">
              Keep your entire team aligned with your codebase. Automatic,
              always fresh.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-8">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="bg-card border border-border rounded-xl p-8 flex flex-col gap-4"
              >
                <Icon
                  type={feature.icon}
                  className="!text-[3rem] text-primary"
                />
                <h3 className={`${headingClass} text-xl font-medium`}>
                  {feature.title}
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  {feature.body}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="py-24 px-4 bg-card">
        <div className="max-w-[900px] mx-auto flex flex-col gap-16">
          <div className="flex flex-col items-center text-center gap-4">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              How it works
            </h2>
            <p className="text-muted-foreground max-w-[500px] leading-relaxed">
              Abstractions connects to your GitHub repositories, deeply indexes
              your codebase, and delivers AI-generated newsletters to your team
              on a schedule you control.
            </p>
          </div>

          {/* Block diagram */}
          <div className="flex flex-col items-center gap-0 w-full">
            {/* Row 1: GitHub Repo */}
            <div className="flex justify-center w-full">
              <div className="bg-background border border-border rounded-xl px-6 py-4 flex flex-col items-center gap-2 w-48 text-center shadow-sm">
                <Icon
                  type="folder_code"
                  className="!text-[2rem] text-primary"
                />
                <span className="font-medium text-sm">GitHub Repository</span>
                <span className="text-xs text-muted-foreground">
                  Public or private
                </span>
              </div>
            </div>

            {/* Arrow down */}
            <div className="flex flex-col items-center gap-0.5 py-1">
              <div className="w-px h-6 bg-border" />
              <div className="text-muted-foreground text-xs">OAuth</div>
              <div className="w-px h-6 bg-border" />
              <div className="text-border">▼</div>
            </div>

            {/* Row 2: Abstractions indexer */}
            <div className="flex flex-col items-center gap-3 w-full">
              <div className="bg-background border-2 border-primary/40 rounded-xl px-6 py-4 flex flex-col items-center gap-2 w-64 text-center shadow-sm">
                <Icon
                  type="manage_search"
                  className="!text-[2rem] text-primary"
                />
                <span className="font-medium text-sm">
                  Abstractions Indexer
                </span>
                <span className="text-xs text-muted-foreground">
                  Semantic & structural analysis
                </span>
              </div>

              {/* Sub-nodes fanning out */}
              <div className="flex items-start justify-center gap-4 w-full">
                {/* left line */}
                <div className="flex flex-col items-end flex-1 max-w-[180px]">
                  <div className="w-full flex justify-end">
                    <div className="flex flex-col items-center gap-0.5">
                      <div className="w-px h-6 bg-border" />
                      <div className="text-border text-xs">▼</div>
                    </div>
                  </div>
                  <div className="bg-background border border-border rounded-xl px-4 py-3 flex flex-col items-center gap-1.5 w-full text-center">
                    <Icon
                      type="data_object"
                      className="!text-[1.5rem] text-muted-foreground"
                    />
                    <span className="font-medium text-xs">Code Chunks</span>
                    <span className="text-xs text-muted-foreground">
                      Files split into meaningful segments
                    </span>
                  </div>
                </div>

                {/* middle */}
                <div className="flex flex-col items-center">
                  <div className="flex flex-col items-center gap-0.5">
                    <div className="w-px h-6 bg-border" />
                    <div className="text-border text-xs">▼</div>
                  </div>
                  <div className="bg-background border border-border rounded-xl px-4 py-3 flex flex-col items-center gap-1.5 w-40 text-center">
                    <Icon
                      type="psychology"
                      className="!text-[1.5rem] text-muted-foreground"
                    />
                    <span className="font-medium text-xs">Embeddings</span>
                    <span className="text-xs text-muted-foreground">
                      Stored in vector DB
                    </span>
                  </div>
                </div>

                {/* right */}
                <div className="flex flex-col items-start flex-1 max-w-[180px]">
                  <div className="w-full flex justify-start">
                    <div className="flex flex-col items-center gap-0.5">
                      <div className="w-px h-6 bg-border" />
                      <div className="text-border text-xs">▼</div>
                    </div>
                  </div>
                  <div className="bg-background border border-border rounded-xl px-4 py-3 flex flex-col items-center gap-1.5 w-full text-center">
                    <Icon
                      type="account_tree"
                      className="!text-[1.5rem] text-muted-foreground"
                    />
                    <span className="font-medium text-xs">Topic Graph</span>
                    <span className="text-xs text-muted-foreground">
                      Concepts & relationships mapped
                    </span>
                  </div>
                </div>
              </div>
            </div>

            {/* Arrow down */}
            <div className="flex flex-col items-center gap-0.5 py-1 mt-2">
              <div className="w-px h-6 bg-border" />
              <div className="text-muted-foreground text-xs">
                Scheduled digest
              </div>
              <div className="w-px h-6 bg-border" />
              <div className="text-border">▼</div>
            </div>

            {/* Row 3: Newsletter generation */}
            <div className="flex justify-center w-full">
              <div className="bg-background border border-border rounded-xl px-6 py-4 flex flex-col items-center gap-2 w-56 text-center shadow-sm">
                <Icon
                  type="auto_awesome"
                  className="!text-[2rem] text-primary"
                />
                <span className="font-medium text-sm">AI Newsletter</span>
                <span className="text-xs text-muted-foreground">
                  Generated by LLM from indexed context
                </span>
              </div>
            </div>

            {/* Arrow down */}
            <div className="flex flex-col items-center gap-0.5 py-1">
              <div className="w-px h-6 bg-border" />
              <div className="text-muted-foreground text-xs">Every Week</div>
              <div className="w-px h-6 bg-border" />
              <div className="text-border">▼</div>
            </div>

            {/* Row 4: Team inboxes */}
            <div className="flex justify-center gap-4 w-full flex-wrap">
              <div className="bg-background border border-border rounded-xl px-5 py-3 flex flex-col items-center gap-1.5 text-center shadow-sm">
                <Icon
                  type="mark_email_read"
                  className="!text-[1.5rem] text-primary"
                />
                <span className="font-medium text-xs">
                  Your Engineering Team
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section id="faq" className="py-24 px-4">
        <div className="max-w-[900px] mx-auto flex flex-col gap-12">
          <div className="flex flex-col items-center text-center gap-4">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              Frequently asked questions
            </h2>
            <p className="text-muted-foreground max-w-[500px] leading-relaxed">
              Everything you need to know about how Abstractions works.
            </p>
          </div>

          <div className="flex flex-col gap-6">
            {[
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
            ].map((faq) => (
              <div
                key={faq.question}
                className="bg-card border border-border rounded-xl p-8 flex flex-col gap-3"
              >
                <h3 className={`${headingClass} text-lg font-medium`}>
                  {faq.question}
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  {faq.answer}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Final CTA */}
      <section className="py-24 px-4">
        <div className="max-w-[900px] mx-auto">
          <div className="bg-card border border-border rounded-xl p-12 flex flex-col items-center text-center gap-6">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              Ready to keep your team in sync?
            </h2>
            <p className="text-muted-foreground max-w-[480px] leading-relaxed">
              Connect your codebase and your first newsletter ships next Monday.
            </p>
            <div className="flex flex-col sm:flex-row items-center gap-3">
              <Link href="/signup">
                <Button size="lg">Sign up free</Button>
              </Link>
              {/* View on GitHub — hidden until open source launch */}
              {/* <Link
                href="https://github.com/nshntarora/abstractions"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="lg" variant="outline">
                  View on GitHub
                </Button>
              </Link> */}
            </div>
            {/* AGPL license text — hidden until open source launch */}
            {/* <p className="text-xs text-muted-foreground">
              AGPL licensed. Audit the code yourself.
            </p> */}
          </div>
        </div>
      </section>

      {/* Footer */}
      <Footer />
    </div>
  );
}
