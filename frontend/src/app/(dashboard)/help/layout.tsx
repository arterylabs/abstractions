"use client";

import { helpNavigation } from "@/config/help-navigation";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function HelpLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const basePath = "/help";

  return (
    <div className="flex gap-8">
      {/* Sidebar */}
      <aside className="hidden md:block w-52 shrink-0">
        <div className="sticky top-4 max-h-[calc(100vh-5rem)] overflow-y-auto pr-2">
          <nav className="flex flex-col gap-0.5">
            {helpNavigation.map((item, idx) =>
              item.type === "heading" ? (
                <p
                  key={idx}
                  className={`${idx === 0 ? "mt-0" : "mt-5"} mb-1.5 text-[0.65rem] font-semibold uppercase tracking-widest text-foreground/35`}
                >
                  {item.title}
                </p>
              ) : (
                <Link
                  key={idx}
                  href={item.slug ? `${basePath}/${item.slug}` : basePath}
                  className={`flex items-center gap-2 rounded-md px-3 py-1.5 text-sm no-underline transition-colors ${
                    pathname ===
                    (item.slug ? `${basePath}/${item.slug}` : basePath)
                      ? "border-l-2 border-primary bg-primary/10 text-primary font-medium"
                      : "text-foreground/60 hover:text-foreground hover:bg-accent/60 border-l-2 border-transparent"
                  }`}
                >
                  {item.title}
                </Link>
              ),
            )}
          </nav>
        </div>
      </aside>

      {/* Content */}
      <div className="flex-1 min-w-0 docs-content max-w-[720px]">
        {children}
      </div>
    </div>
  );
}
