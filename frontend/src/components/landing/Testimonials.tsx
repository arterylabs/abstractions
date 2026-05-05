const testimonials = [
  {
    quote:
      "I was planning my wedding across two cities and the guest management was a nightmare. Harmony changed everything — I could see RSVPs in real time, collect everyone's travel details without a single follow-up call.",
    name: "Meera Krishnan",
    role: "Bride",
    event: "Wedding in Bangalore, 2024",
    emoji: "💍",
    avatarColor: "bg-purple-100",
    avatarText: "MK",
  },
  {
    quote:
      "I threw my best friend's baby shower and used Harmony for the first time. The invite page looked SO good — everyone kept complimenting how beautiful it was. And collecting dietary preferences was so easy!",
    name: "Rhea Sharma",
    role: "Event Host",
    event: "Baby Shower in Mumbai",
    emoji: "🍼",
    avatarColor: "bg-blue-100",
    avatarText: "RS",
  },
  {
    quote:
      "Organised a bachelorette trip for 18 girls across three states. Sent the invite via WhatsApp and had all their travel details and ID proofs collected within 24 hours. Absolute lifesaver.",
    name: "Tanya Oberoi",
    role: "Maid of Honour",
    event: "Bachelorette, Goa 2025",
    emoji: "🥂",
    avatarColor: "bg-pink-100",
    avatarText: "TO",
  },
];

export default function Testimonials() {
  return (
    <section id="testimonials" className="w-full py-16 md:py-24 bg-muted/30">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Real people, real celebrations
          </h2>
          <p className="text-lg text-muted-foreground max-w-lg mx-auto">
            From weddings with 200 guests to birthday dinners with 8—hosts tell
            us the same thing: this made everything easier, and made everything
            feel more personal.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-6">
          {testimonials.map((t) => (
            <div
              key={t.name}
              className="bg-white rounded-2xl p-6 border border-border flex flex-col gap-4 hover:shadow-md transition-shadow"
            >
              {/* Event type badge */}
              <span className="text-2xl">{t.emoji}</span>

              {/* Quote */}
              <blockquote className="text-sm text-foreground leading-relaxed flex-1">
                &ldquo;{t.quote}&rdquo;
              </blockquote>

              {/* Author */}
              <div className="flex items-center gap-3 pt-3 border-t border-border">
                <div
                  className={`w-10 h-10 rounded-full ${t.avatarColor} flex items-center justify-center text-xs font-bold text-foreground flex-shrink-0`}
                >
                  {t.avatarText}
                </div>
                <div>
                  <p className="font-semibold text-sm text-foreground">
                    {t.name}
                  </p>
                  <p className="text-xs text-muted-foreground">{t.event}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
