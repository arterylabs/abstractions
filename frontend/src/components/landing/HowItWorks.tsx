const steps = [
  {
    number: "01",
    emoji: "🎨",
    title: "Tell your story",
    description:
      "Share what matters: the date, the place, the vibe. Pick a design that reflects the feeling you want people to have the moment they see it.",
  },
  {
    number: "02",
    emoji: "📨",
    title: "Invite them in",
    description:
      "Send invites through WhatsApp, SMS, or email—whatever feels natural. No apps to download. Just an invitation that feels like it came from a friend.",
  },
  {
    number: "03",
    emoji: "🎊",
    title: "Stop stressing, start celebrating",
    description:
      "Watch confirmations flow in. Get travel details, dietary needs, everything you need without a single follow-up text. Then put your phone down and actually be present.",
  },
];

export default function HowItWorks() {
  return (
    <section
      id="how-it-works"
      className="w-full py-16 md:py-24 bg-gradient-to-b from-white to-purple-50/50"
    >
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="text-center mb-14">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Planning made simple
          </h2>
          <p className="text-lg text-muted-foreground max-w-lg mx-auto">
            From &ldquo;Save the Date&rdquo; to RSVP confirmed — get everything
            set up in under 10 minutes.
          </p>
        </div>

        {/* Steps */}
        <div className="relative">
          {/* Connecting line (desktop) */}
          <div className="hidden md:block absolute top-12 left-[calc(16.66%+1rem)] right-[calc(16.66%+1rem)] h-0.5 bg-border z-0" />

          <div className="grid md:grid-cols-3 gap-8 relative z-10">
            {steps.map((step) => (
              <div
                key={step.number}
                className="flex flex-col items-center text-center gap-4"
              >
                {/* Step number circle */}
                <div className="w-24 h-24 rounded-full bg-white border-2 border-primary/20 flex flex-col items-center justify-center shadow-sm">
                  <span className="text-3xl">{step.emoji}</span>
                </div>

                <div>
                  <p className="text-xs font-bold tracking-widest text-primary/60 uppercase mb-2">
                    Step {step.number}
                  </p>
                  <h3 className="text-xl font-bold text-foreground mb-2">
                    {step.title}
                  </h3>
                  <p className="text-muted-foreground text-sm leading-relaxed">
                    {step.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
