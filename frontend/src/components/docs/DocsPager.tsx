import Link from "next/link";
import { helpNavigation } from "@/config/help-navigation";
import Icon from "@/components/ui/icon";

// Extract ordered page entries (skip headings)
const pages = helpNavigation.filter((item) => item.type !== "heading") as {
  slug: string;
  title: string;
}[];

function getHref(slug: string, basePath: "/docs" | "/help") {
  return slug ? `${basePath}/${slug}` : basePath;
}

interface DocsPagerProps {
  slug: string;
  basePath?: "/docs" | "/help";
}

export default function DocsPager({
  slug,
  basePath = "/docs",
}: DocsPagerProps) {
  const currentIndex = pages.findIndex((p) => p.slug === slug);
  const prev = currentIndex > 0 ? pages[currentIndex - 1] : null;
  const next = currentIndex < pages.length - 1 ? pages[currentIndex + 1] : null;

  if (!prev && !next) return null;

  return (
    <div className="mt-12 pt-6 border-t border-border grid grid-cols-2 gap-4">
      {prev ? (
        <Link
          href={getHref(prev.slug, basePath)}
          className="group flex flex-col gap-1 p-4 rounded-lg border border-border hover:border-primary/50 hover:bg-accent/50 transition-colors no-underline"
        >
          <span className="text-xs text-foreground/40 uppercase tracking-widest flex items-center gap-1">
            <Icon type="arrow_back" className="text-sm" />
            Previous
          </span>
          <span className="text-sm font-medium text-foreground/80 group-hover:text-primary transition-colors">
            {prev.title}
          </span>
        </Link>
      ) : (
        <div />
      )}
      {next ? (
        <Link
          href={getHref(next.slug, basePath)}
          className="group flex flex-col gap-1 p-4 rounded-lg border border-border hover:border-primary/50 hover:bg-accent/50 transition-colors no-underline text-right"
        >
          <span className="text-xs text-foreground/40 uppercase tracking-widest flex items-center gap-1 justify-end">
            Next
            <Icon type="arrow_forward" className="text-sm" />
          </span>
          <span className="text-sm font-medium text-foreground/80 group-hover:text-primary transition-colors">
            {next.title}
          </span>
        </Link>
      ) : (
        <div />
      )}
    </div>
  );
}
