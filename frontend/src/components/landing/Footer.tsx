import Link from "next/link";
import Logo from "@/components/ui/logo";
import { Separator } from "@/components/ui/separator";

export default function Footer() {
  return (
    <footer className="w-full border-t bg-muted/30">
      <div className="container max-w-[1000px] mx-auto px-4 py-10">
        <div className="flex flex-col md:flex-row justify-between gap-8 mb-8">
          {/* Brand */}
          <div className="flex flex-col gap-3">
            <Link href="/">
              <Logo width="130px" />
            </Link>
            <p className="text-sm text-muted-foreground max-w-xs leading-relaxed">
              Beautiful event invitations for your most personal celebrations.
              Weddings, baby showers, bachelorettes &amp; more.
            </p>
          </div>

          {/* Links */}
          <div className="flex flex-col sm:flex-row gap-8">
            <div>
              <p className="text-sm font-semibold text-foreground mb-3">
                Product
              </p>
              <div className="flex flex-col gap-2">
                <Link
                  href="#occasions"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Occasions
                </Link>
                <Link
                  href="#how-it-works"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  How It Works
                </Link>
                <Link
                  href="#testimonials"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Reviews
                </Link>
              </div>
            </div>
            <div>
              <p className="text-sm font-semibold text-foreground mb-3">
                Account
              </p>
              <div className="flex flex-col gap-2">
                <Link
                  href="/signup"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Sign Up
                </Link>
                <Link
                  href="/login"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Login
                </Link>
              </div>
            </div>
            <div>
              <p className="text-sm font-semibold text-foreground mb-3">
                Legal
              </p>
              <div className="flex flex-col gap-2">
                <Link
                  href="/privacy"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Privacy Policy
                </Link>
                <Link
                  href="/terms"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Terms of Service
                </Link>
              </div>
            </div>
          </div>
        </div>

        <Separator />

        <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-6">
          <p className="text-xs text-muted-foreground">
            © {new Date().getFullYear()} Harmony Invites. All rights reserved.
          </p>
          <p className="text-xs text-muted-foreground">
            Made with 💜 for the moments that matter most
          </p>
        </div>
      </div>
    </footer>
  );
}
