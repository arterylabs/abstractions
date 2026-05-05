const benefits = [
  {
    emoji: "💌",
    headline: "They walk into your world before they walk through the door",
    description:
      "Forget forwarded emails and calendar events. Your guests see your voice, your story, your vibe the moment they open the invite. They're already smiling before they even RSVP yes.",
    highlight: "Beautiful moments that feel deeply personal",
  },
  {
    emoji: "🙌",
    headline: "Everything in one place—so you actually know who's coming",
    description:
      "Stop asking 'Did I hear back from Priya?' five times a day. See who's excited, who's bringing plus-ones, what they need dietary-wise—one glance and you know exactly what's happening. Peace, finally.",
    highlight: "The difference between wondering and knowing",
  },
  {
    emoji: "📲",
    headline: "Meet them where they already are",
    description:
      "Your guests open WhatsApp, tap the invite, and RSVP before they even think about it. No app downloads. No account creation. Just a beautiful moment that feels like it came from you, delivered through whatever they use every day.",
    highlight: "No friction. Just an invitation from a friend.",
  },
];

export default function Benefits() {
  return (
    <section className="w-full py-16 md:py-24 bg-white">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Invitations that do the heavy lifting
          </h2>
          <p className="text-lg text-muted-foreground max-w-xl mx-auto">
            Spend less time managing logistics and more time enjoying the moment
            with the people you love.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-6">
          {benefits.map((benefit) => (
            <div
              key={benefit.headline}
              className="flex flex-col gap-4 p-6 rounded-2xl border border-border hover:border-primary/30 hover:shadow-sm transition-all"
            >
              <span className="text-5xl">{benefit.emoji}</span>
              <div>
                <h3 className="text-xl font-bold text-foreground mb-2">
                  {benefit.headline}
                </h3>
                <p className="text-muted-foreground text-sm leading-relaxed">
                  {benefit.description}
                </p>
              </div>
              <div className="mt-auto pt-3 border-t border-border">
                <p className="text-primary text-sm font-medium">
                  ✓ {benefit.highlight}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
