import Link from "next/link";
import { Button } from "@/components/ui/button";

function InviteCardMockup() {
  return (
    <div className="relative flex items-center justify-center">
      {/* Decorative floating elements */}
      <div className="absolute -top-6 -right-2 text-3xl select-none">🌸</div>
      <div className="absolute top-1/2 -right-8 text-2xl select-none">✨</div>
      <div className="absolute -bottom-4 -left-6 text-3xl select-none">🌿</div>
      <div className="absolute top-8 -left-4 text-2xl select-none">💫</div>

      {/* Main invite card */}
      <div className="bg-white rounded-3xl shadow-2xl p-8 w-72 border border-purple-100 relative z-10">
        <div className="text-center mb-6">
          <div className="text-5xl mb-4">💐</div>
          <p className="text-xs uppercase tracking-[0.2em] text-muted-foreground mb-2 font-medium">
            You&apos;re Invited
          </p>
          <h3 className="text-2xl font-bold text-foreground leading-tight">
            Sarah &amp; James
          </h3>
          <p className="text-primary font-medium text-sm mt-1">
            Wedding Celebration
          </p>
        </div>

        <div className="border-t border-b border-border py-4 mb-6 space-y-2.5">
          <div className="flex items-center gap-2.5 text-sm text-muted-foreground">
            <span className="text-base">📅</span>
            <span>Saturday, June 14, 2025</span>
          </div>
          <div className="flex items-center gap-2.5 text-sm text-muted-foreground">
            <span className="text-base">📍</span>
            <span>The Rosewood Garden, Mumbai</span>
          </div>
          <div className="flex items-center gap-2.5 text-sm text-muted-foreground">
            <span className="text-base">🕕</span>
            <span>6:00 PM onwards</span>
          </div>
        </div>

        <button className="w-full bg-primary text-white rounded-xl py-2.5 text-sm font-semibold hover:bg-primary/90 transition-colors">
          RSVP Now →
        </button>

        <p className="text-center text-xs text-muted-foreground mt-3">
          Hosted on Harmony Invites
        </p>
      </div>

      {/* Small floating RSVP badge */}
      <div className="absolute -bottom-2 right-4 bg-green-50 border border-green-200 rounded-full px-3 py-1.5 flex items-center gap-1.5 shadow-sm z-20">
        <span className="w-2 h-2 rounded-full bg-green-500 inline-block"></span>
        <span className="text-xs text-green-700 font-medium">
          42 guests confirmed
        </span>
      </div>
    </div>
  );
}

export default function Hero() {
  return (
    <section className="w-full py-16 md:py-24 lg:py-32 overflow-hidden">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left: copy */}
          <div className="flex flex-col gap-6">
            <div className="inline-flex items-center gap-2 bg-purple-50 text-primary rounded-full px-4 py-1.5 text-sm font-medium w-fit border border-purple-100">
              <span>💜</span>
              <span>For weddings, baby showers &amp; more</span>
            </div>

            <h1 className="text-4xl md:text-5xl font-bold leading-tight tracking-tight text-foreground">
              Your celebration, <span className="text-primary">done right</span>
            </h1>

            <p className="text-lg text-muted-foreground leading-relaxed">
              Stop juggling WhatsApp threads. Your guests deserve to feel the
              anticipation, not wade through chaos. Send beautiful invitations,
              collect RSVPs, gather details—everything flows like a conversation
              with someone who actually cares about making this special.
            </p>

            <div className="flex flex-col sm:flex-row gap-3 pt-2">
              <Link href="/signup">
                <Button size="lg" className="w-full sm:w-auto text-base px-8">
                  Start planning for free
                </Button>
              </Link>
              <Link href="#how-it-works">
                <Button
                  size="lg"
                  variant="ghost"
                  className="w-full sm:w-auto text-base"
                >
                  See how it works →
                </Button>
              </Link>
            </div>

            <p className="text-sm text-muted-foreground">
              No credit card required · Free to get started
            </p>
          </div>

          {/* Right: invite card mockup */}
          <div className="flex justify-center lg:justify-end">
            <InviteCardMockup />
          </div>
        </div>
      </div>
    </section>
  );
}
