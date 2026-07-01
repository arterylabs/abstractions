import type { Metadata } from "next";
import Script from "next/script";
import { Inter, Geist, Geist_Mono } from "next/font/google";
import "material-symbols/rounded.css";
import "./overrides.css";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";
import { GlobalConfirmationManager } from "@/components/ui/confirmation";
import { ogImages, siteUrl, twitterImages } from "./shared-metadata";

const bodyFont = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

// Crisp geometric sans for headings/display (modern dev-tool look)
const headingFont = Geist({
  variable: "--font-geist",
  subsets: ["latin"],
});

// Monospace for eyebrows, code, and technical accents
const monoFont = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default:
      "Abstractions — Receive weekly newsletters explaining your codebase",
    template: "%s | Abstractions",
  },
  description:
    "Abstractions indexes your GitHub repositories and delivers weekly AI-generated newsletters that explain your codebase, keeping your entire team — humans and agents alike — in sync as it evolves.",
  keywords: [
    "codebase knowledge",
    "developer newsletters",
    "AI code summaries",
    "context debt",
    "engineering team",
    // "open source", // hidden until open source launch
  ],
  authors: [{ name: "Abstractions" }],
  robots: {
    index: true,
    follow: true,
    googleBot: { index: true, follow: true },
  },
  openGraph: {
    type: "website",
    siteName: "Abstractions",
    title: "Abstractions — Receive weekly newsletters explaining your codebase",
    description:
      "Weekly AI-generated newsletters that explain your codebase, keeping your entire engineering team in sync as it evolves.",
    images: [...ogImages],
  },
  twitter: {
    card: "summary_large_image",
    title: "Abstractions — Receive weekly newsletters explaining your codebase",
    description:
      "Weekly AI-generated newsletters that explain your codebase, keeping your entire engineering team in sync as it evolves.",
    images: twitterImages,
  },
};

const organizationSchema = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: "Abstractions",
  url: siteUrl,
  sameAs: ["https://github.com/nshntarora/abstractions"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <head>
        <Script
          defer
          data-domain="useabstractions.com"
          src="/p/script.js"
          strategy="afterInteractive"
        />
        <Script id="plausible-init" strategy="afterInteractive">
          {`
            window.plausible = window.plausible || function() { (window.plausible.q = window.plausible.q || []).push(arguments) };
            window.plausible.init = window.plausible.init || function(i) { window.plausible.o = i || {} };
            window.plausible.init({ endpoint: "/p/event" });
          `}
        </Script>
      </head>
      <body
        className={`${bodyFont.variable} ${headingFont.variable} ${monoFont.variable} antialiased`}
      >
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify(organizationSchema),
          }}
        />
        {children}
        <Toaster />
        <GlobalConfirmationManager />
        <div id="portal-root" />
      </body>
    </html>
  );
}
