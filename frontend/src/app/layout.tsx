import type { Metadata, Viewport } from "next";
import Script from "next/script";
import {
  Instrument_Serif,
  IBM_Plex_Sans,
  IBM_Plex_Mono,
} from "next/font/google";
import "material-symbols/rounded.css";
import "./overrides.css";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";
import { GlobalConfirmationManager } from "@/components/ui/confirmation";
import { ogImages, siteUrl, twitterImages } from "./shared-metadata";

// Matches --background in globals.css. Applied inline so the canvas is dark
// before CSS loads (avoids a white flash on first paint).
const PAGE_BACKGROUND = "#08090a";

/* "Night press" typography — one editorial voice across the marketing site,
   the app, and the docs. See the `.press-*` utilities in globals.css. */

// Body / UI text — technical, but warmer than a neutral grotesque.
const bodyFont = IBM_Plex_Sans({
  variable: "--font-plex-sans",
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  display: "swap",
});

// High-contrast editorial serif for display headlines and section titles.
const headingFont = Instrument_Serif({
  variable: "--font-instrument",
  subsets: ["latin"],
  weight: "400",
  style: ["normal", "italic"],
  display: "swap",
});

// Datelines, eyebrows, folios, file paths, identifiers, timestamps.
const monoFont = IBM_Plex_Mono({
  variable: "--font-plex-mono",
  subsets: ["latin"],
  weight: ["400", "500"],
  display: "swap",
});

export const viewport: Viewport = {
  colorScheme: "dark",
  themeColor: PAGE_BACKGROUND,
};

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
    "open source",
    "AGPL",
    "self-host",
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
  sameAs: ["https://github.com/arterylabs/abstractions"],
};

const analyticsEnabled = process.env.NEXT_PUBLIC_ENABLE_ANALYTICS === "true";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    // The next/font variables must sit on <html>, not <body>: globals.css
    // declares --font-body / --font-heading on :root, and the var() references
    // inside those values are resolved against :root itself. On <body> they
    // would be invisible there and the stacks would silently lose their faces.
    <html
      lang="en"
      className={`dark ${bodyFont.variable} ${headingFont.variable} ${monoFont.variable}`}
      style={{ backgroundColor: PAGE_BACKGROUND, colorScheme: "dark" }}
    >
      <head>
        {/* Critical inline style so first paint is dark even before CSS bundles load */}
        <style
          dangerouslySetInnerHTML={{
            __html: `html,body{background-color:${PAGE_BACKGROUND};color-scheme:dark}`,
          }}
        />
        {analyticsEnabled && (
          <>
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
          </>
        )}
      </head>
      <body
        className="antialiased"
        style={{ backgroundColor: PAGE_BACKGROUND }}
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
