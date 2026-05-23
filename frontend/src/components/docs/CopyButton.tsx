"use client";

import { useState } from "react";
import Icon from "@/components/ui/icon";

export default function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <button
      onClick={handleCopy}
      aria-label="Copy code"
      className="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity bg-accent hover:bg-accent-light border border-border rounded p-1 text-foreground/60 hover:text-foreground"
    >
      <Icon type={copied ? "check" : "content_copy"} className="text-sm" />
    </button>
  );
}
