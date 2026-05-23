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

  return (
    <div className="grid grid-cols-5 gap-6">
      <div className="col-span-1 flex flex-col gap-4 border-r">
        {helpNavigation.map((item, idx) =>
          item.type === "heading" ? (
            <div
              key={idx}
              className="font-medium uppercase text-sm opacity-60 mt-6"
            >
              {item.title}
            </div>
          ) : (
            <div key={idx} className="mt-2">
              <Link
                href={item.slug ? `/help/${item.slug}` : "/help"}
                className={`hover:text-primary ${
                  pathname === (item.slug ? `/help/${item.slug}` : "/help")
                    ? "text-primary border-b border-primary"
                    : ""
                }`}
              >
                {item.title}
              </Link>
            </div>
          ),
        )}
      </div>
      <div className="col-span-4 docs-content">{children}</div>
    </div>
  );
}
