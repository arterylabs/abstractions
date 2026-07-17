"use client";

import { cn } from "@/lib/utils";
import DOMPurify from "dompurify";

interface RichTextContentProps {
  html?: string | null;
  className?: string;
}

/**
 * Renders stored TipTap HTML content safely.
 * DOMPurify strips any script tags or event handlers before the HTML is
 * injected into the DOM, preventing stored-XSS attacks.
 * Use this component on public/read-only pages to display rich text.
 */
export function RichTextContent({ html, className }: RichTextContentProps) {
  if (!html) return null;

  const sanitizedHtml = DOMPurify.sanitize(html);

  return (
    <div
      className={cn(
        "prose prose-sm max-w-none text-foreground",
        "[&_h1]:text-2xl [&_h1]:font-bold [&_h1]:mb-2",
        "[&_h2]:text-xl [&_h2]:font-semibold [&_h2]:mb-2",
        "[&_h3]:text-lg [&_h3]:font-semibold [&_h3]:mb-1",
        "[&_p]:mb-2 [&_ul]:list-disc [&_ul]:pl-5 [&_ol]:list-decimal [&_ol]:pl-5",
        "[&_blockquote]:border-l-4 [&_blockquote]:border-primary/40 [&_blockquote]:pl-4 [&_blockquote]:italic [&_blockquote]:text-muted-foreground [&_blockquote]:my-3",
        className,
      )}
      dangerouslySetInnerHTML={{ __html: sanitizedHtml }}
    />
  );
}
