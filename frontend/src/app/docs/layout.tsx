"use client";

import Footer from "@/components/landing/Footer";
import LandingPageNavbar from "@/components/landing/Navbar";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function Docs({ children }: { children: React.ReactNode }) {
  const sidebarLinks = [
    {
      type: "heading",
      title: "Getting Started",
    },
    {
      link: "/docs",
      title: "Introduction",
    },
    {
      link: "/docs/installation",
      title: "Installation",
    },
    {
      type: "heading",
      title: "Configure",
    },
    {
      link: "/docs/google",
      title: "Google Search",
    },
    {
      link: "/docs/openai",
      title: "OpenAI",
    },
    {
      type: "heading",
      title: "Usage",
    },
    {
      link: "/docs/import",
      title: "Importing Data",
    },
    {
      link: "/docs/enrichment",
      title: "Enriching Records",
    },
    {
      link: "/docs/export",
      title: "Exporting Data",
    },
  ];

  const pathname = usePathname();

  return (
    <div className="flex flex-col">
      <LandingPageNavbar />
      <div className="grid grid-cols-5 max-w-[1000px] mx-auto container py-12 gap-6">
        {/* <div className="col-span-5 font-medium text-3xl pb-4 mb-4">
          Documentation
        </div> */}
        <div className="col-span-1 min-h-screen h-full flex flex-col gap-4 border-r">
          {sidebarLinks.map((link, idx) =>
            link.type === "heading" ? (
              <div
                key={idx}
                className="font-medium uppercase text-sm opacity-60 mt-6"
              >
                {link.title}
              </div>
            ) : (
              <div key={idx} className="mt-2">
                <Link
                  href={link.link || "/docs"}
                  className={`hover:text-primary ${
                    pathname === link.link
                      ? "text-primary border-b border-primary"
                      : ""
                  }`}
                >
                  {link.title}
                </Link>
              </div>
            ),
          )}
        </div>
        <div className="col-span-4 min-h-screen h-full docs-content">
          {children}
        </div>
      </div>
      <Footer />
    </div>
  );
}
