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

export const metadata: Metadata = {
  title: "Abstractions",
  description: "Enrich your leads using AI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${bodyFont.variable} ${headingFont.variable} antialiased`}
      >
        {children}
        <Toaster />
        <GlobalConfirmationManager />
        <div id="portal-root" />
      </body>
    </html>
  );
}
