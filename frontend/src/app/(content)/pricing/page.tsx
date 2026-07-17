import type { Metadata } from "next";
import { ogImages, siteUrl } from "../../shared-metadata";

const GITHUB_URL = "https://github.com/arterylabs/abstractions";

export const metadata: Metadata = {
  title: "Pricing — Free to use",
  description:
    "Abstractions is free to use. The hosted product runs the same open-source code as the GitHub repository. Bring your own model keys — we only cover hosting.",
  alternates: {
    canonical: `${siteUrl}/pricing`,
  },
  openGraph: {
    type: "website",
    url: `${siteUrl}/pricing`,
    title: "Pricing — Free to use | Abstractions",
    description:
      "The hosted product is free. Same codebase as GitHub. Bring your own model keys.",
    images: ogImages,
  },
};

export default function PricingPage() {
  return (
    <main className="flex flex-1 flex-col bg-background text-foreground">
      <section className="px-4 py-24">
        <div className="mx-auto flex max-w-[760px] flex-col gap-8">
          <div className="flex flex-col items-center gap-4 text-center">
            <span className="font-mono text-xs uppercase tracking-[0.2em] text-primary">
              Pricing
            </span>
            <h1 className="text-3xl font-semibold md:text-4xl">
              Bring your own key. Use for free.
            </h1>
          </div>

          <div className="flex flex-col gap-5 rounded-2xl border border-border bg-card p-8">
            <p className="leading-relaxed text-muted-foreground">
              Since the whole app is bring your own keys (which are encrypted),
              the only cost I have is hosting. It currently costs &lt;$50 to
              host this whole app (two instances, one for the database, one for
              the api, frontend, and workers).
            </p>
            <p className="leading-relaxed text-muted-foreground">
              What you see on the hosted site is the same code published on{" "}
              <a
                href={GITHUB_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="font-medium text-primary underline-offset-2 hover:underline"
              >
                GitHub
              </a>
              . Inspect it, fork it, or self-host it. Nothing important lives
              only behind our domain.
            </p>

            <p className="leading-relaxed text-muted-foreground">
              If there&apos;s new functionality and increased usage that needs
              more resources from my end, I&apos;d come up with pricing (with
              enough notice). For now, use it for free.
            </p>
          </div>
        </div>
      </section>
    </main>
  );
}
