import Link from "next/link";
import Icon from "@/components/ui/icon";
import { helpNavigation } from "@/config/help-navigation";

interface DocsBreadcrumbsProps {
  slug: string;
  basePath?: "/docs" | "/help";
}

export default function DocsBreadcrumbs({
  slug,
  basePath = "/docs",
}: DocsBreadcrumbsProps) {
  // Find section heading for the current page
  let sectionTitle = "";
  let pageTitle = "";

  let currentHeading = "";
  for (const item of helpNavigation) {
    if (item.type === "heading") {
      currentHeading = item.title;
    } else if (item.slug === slug) {
      sectionTitle = currentHeading;
      pageTitle = item.title;
      break;
    }
  }

  // Root page ("Introduction") has empty slug
  if (!slug) {
    pageTitle =
      (
        helpNavigation.find((i) => i.type !== "heading" && i.slug === "") as
          | { title: string }
          | undefined
      )?.title ?? "Introduction";
    sectionTitle =
      helpNavigation[0]?.type === "heading" ? helpNavigation[0].title : "";
  }

  return (
    <nav
      aria-label="Breadcrumb"
      className="flex items-center gap-1.5 text-xs text-foreground/40 mb-6"
    >
      <Link
        href={basePath}
        className="hover:text-foreground/70 transition-colors no-underline"
      >
        Docs
      </Link>
      {sectionTitle && (
        <>
          <Icon type="chevron_right" className="text-xs" />
          <span>{sectionTitle}</span>
        </>
      )}
      {pageTitle && (
        <>
          <Icon type="chevron_right" className="text-xs" />
          <span className="text-foreground/70">{pageTitle}</span>
        </>
      )}
    </nav>
  );
}
