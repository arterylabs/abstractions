import Link from "next/link";
import Logo from "@/components/ui/logo";

const footerLinks = [
  { href: "/docs", label: "Docs" },
  { href: "/privacy", label: "Privacy" },
  { href: "/terms", label: "Terms" },
  { href: "/cookies", label: "Cookies" },
  { href: "/login", label: "Login" },
];

export default function Footer() {
  return (
    <footer className="border-t border-border p-8">
      <div className="max-w-[1200px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
        <div className="flex items-center gap-2">
          <Logo className="h-5 w-auto opacity-70" />
          <span>© {new Date().getFullYear()} Artery Ventures, LLP</span>
        </div>
        <div className="flex items-center gap-6">
          {footerLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="hover:text-foreground transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </div>
      </div>
    </footer>
  );
}
