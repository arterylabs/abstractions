import Link from "next/link";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";
import LandingPageNavbar from "@/components/landing/Navbar";

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
    icon: "smart_toy",
    title: "AI agents operate in a vacuum",
    body: "Without codebase context, agents make decisions that silently undo what you just shipped last sprint.",
  },
];

const features = [
  {
    icon: "code",
    title: "Open source & self-hostable",
    body: "Deploy with Docker in minutes. No vendor lock-in, no usage limits — the entire codebase is on GitHub.",
  },
  {
    icon: "manage_search",
    title: "Deep semantic indexing",
    body: "Abstractions analyzes structure, patterns, and relationships across your entire codebase — not just recent diffs.",
  },
  {
    icon: "mark_email_read",
    title: "AI-generated newsletters",
    body: "Every Monday (or whenever you configure) a clean, readable summary lands in every inbox on your team.",
  },
  {
    icon: "lock",
    title: "Full control over your data",
    body: "Run on your own infrastructure. Your source code and embeddings never leave your environment.",
  },
];

const steps = [
  {
    number: "01",
    title: "Star & clone the repo",
    body: "No account required. Grab the source from GitHub and deploy it wherever you run your stack.",
  },
  {
    number: "02",
    title: "Connect your GitHub repo",
    body: "Public or private. OAuth takes 30 seconds and we request only the minimum permissions needed.",
  },
  {
    number: "03",
    title: "Abstractions indexes your codebase",
    body: "Deep semantic and structural analysis runs automatically in the background.",
  },
  {
    number: "04",
    title: "Receive your newsletter",
    body: "Every Monday (or the day you configure) your team gets a clean, professional digest.",
  },
];

export default function LandingPage() {
  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground">
      <LandingPageNavbar />

      {/* Hero */}
      <section className="py-24 px-4">
        <div className="max-w-[900px] mx-auto flex flex-col items-center text-center gap-6">
          <div className="inline-flex items-center gap-2 border border-primary/40 text-primary text-sm px-3 py-1 rounded-full">
            <span className="size-1.5 rounded-full bg-primary inline-block" />
            Open Source · AGPL License
          </div>

          <h1
            className={`${headingClass} text-2xl md:text-3xl lg:text-4xl font-medium leading-tight`}
          >
            Never forget how your code works <br />
            when you build with agentic speed
          </h1>

          <p className="text-lg text-muted-foreground max-w-[600px] leading-relaxed">
            Receive weekly newsletters that explain concepts in your code
          </p>

          <div className="flex flex-col sm:flex-row items-center gap-3 mt-2">
            <Link href="/signup">
              <Button size="lg">Sign up</Button>
            </Link>
            <Link
              href="https://github.com/nshntarora/abstractions"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Button size="lg" variant="outline">
                View on GitHub
              </Button>
            </Link>
          </div>

          <p className="text-sm text-muted-foreground">
            AGPL licensed. Audit the code yourself.
          </p>

          {/* Screenshot placeholder */}
          <div className="w-full mt-8 bg-card border border-border rounded-xl aspect-video flex items-center justify-center">
            <p className="text-muted-foreground text-sm">
              App screenshot coming soon
            </p>
          </div>
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
              The faster you ship, the harder it becomes for everyone — humans
              and agents alike — to understand how everything fits together.
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
              Everything you need to keep your team in sync
            </h2>
            <p className="text-muted-foreground max-w-[500px] leading-relaxed">
              Automatic Newsletters works in the background so your team never
              has to.
            </p>
          </div>

          <div className="flex flex-col gap-16">
            {features.map((feature, i) => (
              <div
                key={feature.title}
                className={`grid grid-cols-1 md:grid-cols-2 gap-8 items-center ${
                  i % 2 === 1 ? "md:[direction:rtl]" : ""
                }`}
              >
                <div
                  className={`bg-accent border border-border rounded-xl aspect-video flex items-center justify-center ${
                    i % 2 === 1 ? "md:[direction:ltr]" : ""
                  }`}
                >
                  <Icon
                    type={feature.icon}
                    className="text-[4rem] text-primary opacity-60"
                  />
                </div>
                <div
                  className={`flex flex-col gap-3 ${
                    i % 2 === 1 ? "md:[direction:ltr]" : ""
                  }`}
                >
                  <h3 className={`${headingClass} text-xl font-medium`}>
                    {feature.title}
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    {feature.body}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="py-24 px-4 bg-card">
        <div className="max-w-[900px] mx-auto flex flex-col gap-12">
          <div className="flex flex-col items-center text-center gap-4">
            <h2 className={`${headingClass} text-3xl font-medium`}>
              Up and running in minutes
            </h2>
            <p className="text-muted-foreground max-w-[500px] leading-relaxed">
              Self-host with Docker or run it locally. No third-party accounts
              required.
            </p>
          </div>

          <div className="flex flex-col gap-0 max-w-[600px] mx-auto w-full">
            {steps.map((step, i) => (
              <div key={step.number} className="flex gap-6">
                <div className="flex flex-col items-center">
                  <div className="flex items-center justify-center size-10 rounded-full border border-primary/40 bg-primary/10 shrink-0">
                    <span
                      className={`${headingClass} text-sm font-medium text-primary`}
                    >
                      {step.number}
                    </span>
                  </div>
                  {i < steps.length - 1 && (
                    <div className="w-px flex-1 bg-border mt-2 mb-2" />
                  )}
                </div>
                <div
                  className={`flex flex-col gap-1 ${i < steps.length - 1 ? "pb-8" : ""}`}
                >
                  <h3 className="font-medium text-base">{step.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    {step.body}
                  </p>
                </div>
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
              Clone the repo, connect your codebase, and your first newsletter
              ships next Monday.
            </p>
            <div className="flex flex-col sm:flex-row items-center gap-3">
              <Link href="/signup">
                <Button size="lg">Sign up free</Button>
              </Link>
              <Link
                href="https://github.com/nshntarora/abstractions"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="lg" variant="outline">
                  View on GitHub
                </Button>
              </Link>
            </div>
            <p className="text-xs text-muted-foreground">
              AGPL licensed. Audit the code yourself.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border py-8 px-4">
        <div className="max-w-[900px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
          <span>© {new Date().getFullYear()} Abstractions</span>
          <div className="flex items-center gap-6">
            <Link
              href="https://github.com/nshntarora/abstractions"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </Link>
            <Link
              href="/privacy"
              className="hover:text-foreground transition-colors"
            >
              Privacy
            </Link>
            <Link
              href="/terms"
              className="hover:text-foreground transition-colors"
            >
              Terms
            </Link>
            <Link
              href="/login"
              className="hover:text-foreground transition-colors"
            >
              Login
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
