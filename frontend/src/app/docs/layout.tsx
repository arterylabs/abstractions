"use client";

import LandingPageNavbar from "@/components/landing/Navbar";
import { helpNavigation } from "@/config/help-navigation";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  return (
    <div className="flex flex-col">
      <LandingPageNavbar />
      <div className="grid grid-cols-5 max-w-[1000px] mx-auto container py-12 gap-6">
        <div className="col-span-1 min-h-screen h-full flex flex-col gap-4 border-r">
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
                  href={item.slug ? `/docs/${item.slug}` : "/docs"}
                  className={`hover:text-primary ${
                    pathname === (item.slug ? `/docs/${item.slug}` : "/docs")
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
        <div className="col-span-4 min-h-screen h-full docs-content">
          {children}
        </div>
      </div>
    </div>
  );
}
