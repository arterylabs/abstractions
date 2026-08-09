import { helpNavigation } from "@/config/help-navigation";
import Icon from "@/components/ui/icon";
import { Button } from "@/components/ui/button";

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
        <Button
          href={getHref(prev.slug, basePath)}
          variant="outline"
          asChild
          className="group h-auto flex flex-col items-start gap-1 p-4 rounded-lg no-underline"
        >
          <span className="press-label text-muted-foreground/70 flex items-center gap-1">
            <Icon type="arrow_back" className="text-sm" />
            Previous
          </span>
          <span className="text-sm font-medium text-foreground/80 group-hover:text-primary transition-colors">
            {prev.title}
          </span>
        </Button>
      ) : (
        <div />
      )}
      {next ? (
        <Button
          href={getHref(next.slug, basePath)}
          variant="outline"
          asChild
          className="group h-auto flex flex-col items-end gap-1 p-4 rounded-lg no-underline"
        >
          <span className="press-label text-muted-foreground/70 flex items-center gap-1 justify-end">
            Next
            <Icon type="arrow_forward" className="text-sm" />
          </span>
          <span className="text-sm font-medium text-foreground/80 group-hover:text-primary transition-colors">
            {next.title}
          </span>
        </Button>
      ) : (
        <div />
      )}
    </div>
  );
}
