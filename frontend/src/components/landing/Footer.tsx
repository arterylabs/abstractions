import Link from "next/link";
import Logo from "@/components/ui/logo";

const footerLinks = [
  { href: "/docs", label: "Docs" },
  { href: "/pricing", label: "Pricing" },
  { href: "/privacy", label: "Privacy" },
  { href: "/terms", label: "Terms" },
  { href: "/cookies", label: "Cookies" },
  { href: "/login", label: "Login" },
  {
    href: "https://github.com/arterylabs/abstractions",
    label: "GitHub",
    external: true,
  },
];

export default function Footer() {
  return (
    <footer className="border-t border-border px-6 py-12">
      <div className="mx-auto flex max-w-[1200px] flex-col gap-10">
        <div className="flex flex-col items-start justify-between gap-8 md:flex-row md:items-center">
          <div className="flex max-w-sm flex-col items-start justify-start gap-3">
            <Logo className="h-6 w-auto" />
            <p className="text-sm leading-relaxed text-muted-foreground">
              Open-source AI newsletters that keep your engineering team in sync
              with how your codebase evolves.
            </p>
          </div>
          <nav className="flex flex-wrap items-center gap-x-6 gap-y-3 text-sm text-muted-foreground">
            {footerLinks.map((link) =>
              "external" in link && link.external ? (
                <a
                  key={link.href}
                  href={link.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="!text-muted-foreground transition-colors hover:!text-foreground"
                >
                  {link.label}
                </a>
              ) : (
                <Link
                  key={link.href}
                  href={link.href}
                  className="!text-muted-foreground transition-colors hover:!text-foreground"
                >
                  {link.label}
                </Link>
              ),
            )}
          </nav>
        </div>

        <div className="flex flex-col items-center justify-between gap-3 border-t border-border pt-6 text-xs text-muted-foreground sm:flex-row">
          <span>
            © {new Date().getFullYear()} Artery Ventures, LLP ·{" "}
            <a
              href="https://github.com/arterylabs/abstractions/blob/master/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
              className="!text-muted-foreground underline-offset-2 transition-colors hover:!text-foreground hover:underline"
            >
              AGPL-3.0
            </a>
          </span>
          <span className="font-mono">
            Open source · Built for engineering teams
          </span>
        </div>
      </div>
    </footer>
  );
}
