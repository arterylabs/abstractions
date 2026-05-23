import type { Metadata } from "next";
import { Inter, Fira_Code } from "next/font/google";
import "material-symbols/rounded.css";
import "./overrides.css";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";
import { GlobalConfirmationManager } from "@/components/ui/confirmation";

const bodyFont = Inter({
  variable: "--font-body",
  subsets: ["latin"],
});

const headingFont = Fira_Code({
  variable: "--font-heading",
  subsets: ["latin"],
});

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL ?? "https://abstractions.dev";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "Abstractions — Keep your codebase knowledge fresh",
    template: "%s | Abstractions",
  },
  description:
    "Abstractions indexes your GitHub repositories and delivers weekly AI-generated newsletters that keep your entire team — humans and agents alike — in sync with how your codebase evolves.",
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
    title: "Abstractions — Keep your codebase knowledge fresh",
    description:
      "Weekly AI-generated newsletters that keep your entire engineering team in sync with how your codebase evolves.",
    images: [
      {
        url: "/opengraph-image",
        width: 1200,
        height: 630,
        alt: "Abstractions — Keep your codebase knowledge fresh",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Abstractions — Keep your codebase knowledge fresh",
    description:
      "Weekly AI-generated newsletters that keep your entire engineering team in sync with how your codebase evolves.",
    images: ["/opengraph-image"],
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
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin="anonymous"
        />
      </head>
      <body
        className={`${bodyFont.variable} ${headingFont.variable} antialiased`}
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
