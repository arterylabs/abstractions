const occasions = [
  {
    emoji: "💍",
    label: "Weddings",
    description:
      "From two people to a tribe—celebrate the ones who shaped your love story",
    gradient: "from-purple-50 to-pink-50",
    border: "border-purple-100",
  },
  {
    emoji: "🍼",
    label: "Baby Showers",
    description:
      "Welcome your little one with the people who'll help raise them",
    gradient: "from-blue-50 to-cyan-50",
    border: "border-blue-100",
  },
  {
    emoji: "🥂",
    label: "Bachelorettes",
    description: "One last night before everything changes—make it legendary",
    gradient: "from-pink-50 to-rose-50",
    border: "border-pink-100",
  },
  {
    emoji: "🎂",
    label: "Birthdays",
    description: "Milestone birthdays deserve milestone invites",
    gradient: "from-yellow-50 to-orange-50",
    border: "border-yellow-100",
  },
  {
    emoji: "💝",
    label: "Anniversaries",
    description: "Celebrate the years of love",
    gradient: "from-red-50 to-pink-50",
    border: "border-red-100",
  },
  {
    emoji: "🎉",
    label: "Get-togethers",
    description: "Any reason to bring your people together",
    gradient: "from-green-50 to-teal-50",
    border: "border-green-100",
  },
];

export default function Occasions() {
  return (
    <section id="occasions" className="w-full py-16 md:py-24 bg-muted/40">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Made for every celebration
          </h2>
          <p className="text-lg text-muted-foreground max-w-xl mx-auto">
            Whether it&apos;s a once-in-a-lifetime wedding or a cosy birthday
            dinner — Harmony makes it feel special.
          </p>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {occasions.map((occasion) => (
            <div
              key={occasion.label}
              className={`bg-gradient-to-br ${occasion.gradient} border ${occasion.border} rounded-2xl p-6 flex flex-col gap-3 hover:shadow-md transition-shadow`}
            >
              <span className="text-4xl">{occasion.emoji}</span>
              <div>
                <h3 className="font-bold text-foreground text-lg">
                  {occasion.label}
                </h3>
                <p className="text-sm text-muted-foreground mt-1">
                  {occasion.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
