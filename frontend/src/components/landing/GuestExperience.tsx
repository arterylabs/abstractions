function PhoneMockup() {
  return (
    <div className="relative flex justify-center items-center">
      {/* Phone frame */}
      <div className="w-56 bg-white rounded-[2.5rem] shadow-2xl border-4 border-gray-800 overflow-hidden">
        {/* Phone notch */}
        <div className="bg-gray-800 h-6 flex items-center justify-center">
          <div className="w-16 h-1.5 bg-gray-600 rounded-full" />
        </div>

        {/* Phone screen content */}
        <div className="bg-purple-50 min-h-[400px] p-4">
          {/* App bar */}
          <div className="flex items-center justify-between mb-4">
            <span className="text-xs font-bold text-primary">Harmony</span>
            <span className="text-xs text-muted-foreground">9:41 AM</span>
          </div>

          {/* Invite header */}
          <div className="bg-gradient-to-b from-purple-600 to-purple-800 rounded-xl p-4 text-white text-center mb-3">
            <div className="text-2xl mb-2">🌸</div>
            <p className="text-xs font-medium uppercase tracking-wider opacity-80">
              You&apos;re Invited
            </p>
            <h4 className="font-bold text-sm mt-0.5">
              Priya&apos;s Baby Shower
            </h4>
          </div>

          {/* Details */}
          <div className="bg-white rounded-xl p-3 space-y-2 mb-3 shadow-sm">
            <div className="flex items-center gap-2">
              <span className="text-xs">📅</span>
              <span className="text-xs text-muted-foreground">
                Sat, Aug 2 · 3:00 PM
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs">📍</span>
              <span className="text-xs text-muted-foreground">
                The Grand Ballroom
              </span>
            </div>
          </div>

          {/* RSVP button */}
          <button className="w-full bg-primary text-white rounded-xl py-2 text-xs font-semibold">
            RSVP Now
          </button>

          <p className="text-center text-[10px] text-muted-foreground mt-2">
            No app download needed
          </p>
        </div>

        {/* Home indicator */}
        <div className="bg-white h-5 flex items-center justify-center border-t border-gray-100">
          <div className="w-20 h-1 bg-gray-300 rounded-full" />
        </div>
      </div>

      {/* Floating notification */}
      <div className="absolute -right-4 top-16 bg-white rounded-2xl shadow-lg p-3 border border-border flex items-center gap-2 w-44">
        <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center flex-shrink-0">
          <span className="text-sm">✅</span>
        </div>
        <div>
          <p className="text-xs font-semibold text-foreground">
            Anjali RSVP&apos;d
          </p>
          <p className="text-[10px] text-muted-foreground">
            Attending + 1 guest
          </p>
        </div>
      </div>

      {/* Floating detail collected badge */}
      <div className="absolute -left-6 bottom-16 bg-white rounded-2xl shadow-lg p-3 border border-border w-40">
        <p className="text-[10px] font-semibold text-foreground mb-1">
          Details collected
        </p>
        <div className="space-y-0.5">
          <p className="text-[10px] text-muted-foreground flex items-center gap-1">
            <span className="text-green-500">✓</span> Flight details
          </p>
          <p className="text-[10px] text-muted-foreground flex items-center gap-1">
            <span className="text-green-500">✓</span> ID proof
          </p>
          <p className="text-[10px] text-muted-foreground flex items-center gap-1">
            <span className="text-green-500">✓</span> Dietary needs
          </p>
        </div>
      </div>
    </div>
  );
}

export default function GuestExperience() {
  return (
    <section className="w-full py-16 md:py-24 bg-white">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left: phone mockup */}
          <div className="flex justify-center order-2 lg:order-1">
            <PhoneMockup />
          </div>

          {/* Right: copy */}
          <div className="order-1 lg:order-2 flex flex-col gap-6">
            <div className="inline-flex items-center gap-2 bg-green-50 text-green-700 rounded-full px-4 py-1.5 text-sm font-medium w-fit border border-green-100">
              <span>📲</span>
              <span>What your guests experience</span>
            </div>

            <h2 className="text-3xl md:text-4xl font-bold text-foreground leading-tight">
              A beautiful moment—{" "}
              <span className="text-primary">
                not just another link to click
              </span>
            </h2>

            <p className="text-lg text-muted-foreground leading-relaxed">
              Your guests open something beautiful. No app to download, no
              account to create—just a moment that makes them excited to say yes
              before they even finish reading it.
            </p>

            <ul className="space-y-3">
              {[
                "Works on any phone, instantly",
                "RSVP with a single tap",
                "Collect travel details & ID proofs seamlessly",
                "Guests get reminders as the date approaches",
              ].map((point) => (
                <li
                  key={point}
                  className="flex items-center gap-3 text-foreground"
                >
                  <span className="w-5 h-5 rounded-full bg-primary/10 text-primary flex items-center justify-center flex-shrink-0 text-xs font-bold">
                    ✓
                  </span>
                  <span className="text-sm">{point}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}
