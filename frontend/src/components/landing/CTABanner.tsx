import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function CTABanner() {
  return (
    <section className="w-full py-16 md:py-24">
      <div className="container max-w-[1000px] mx-auto px-4">
        <div className="bg-primary rounded-3xl px-8 py-14 md:px-16 text-center relative overflow-hidden">
          {/* Decorative background emojis */}
          <div className="absolute top-4 left-6 text-4xl opacity-20 select-none">
            💍
          </div>
          <div className="absolute top-4 right-8 text-4xl opacity-20 select-none">
            🌸
          </div>
          <div className="absolute bottom-4 left-12 text-4xl opacity-20 select-none">
            🥂
          </div>
          <div className="absolute bottom-4 right-6 text-4xl opacity-20 select-none">
            🎉
          </div>
          <div className="absolute top-1/2 -translate-y-1/2 left-4 text-3xl opacity-10 select-none">
            ✨
          </div>
          <div className="absolute top-1/2 -translate-y-1/2 right-4 text-3xl opacity-10 select-none">
            💫
          </div>

          <div className="relative z-10">
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-4 leading-tight">
              Ready to make your guests feel special?
            </h2>
            <p className="text-purple-100 text-lg mb-8 max-w-lg mx-auto leading-relaxed">
              Create your event page in minutes. Free to get started — no credit
              card needed.
            </p>
            <Link href="/signup">
              <Button
                size="lg"
                className="bg-white text-primary hover:bg-purple-50 text-base px-10 font-semibold shadow-lg"
              >
                Create your event for free →
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}
