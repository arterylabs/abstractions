import Link from "next/link";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
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
          { title: string } | undefined
      )?.title ?? "Introduction";
    sectionTitle =
      helpNavigation[0]?.type === "heading" ? helpNavigation[0].title : "";
  }

  return (
    <Breadcrumb className="mb-6 text-xs">
      <BreadcrumbList className="text-xs flex items-center gap-1">
        <BreadcrumbItem>
          <BreadcrumbLink asChild>
            <Link href={basePath}>Docs</Link>
          </BreadcrumbLink>
        </BreadcrumbItem>
        {sectionTitle && (
          <>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <span>{sectionTitle}</span>
            </BreadcrumbItem>
          </>
        )}
        {pageTitle && (
          <>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{pageTitle}</BreadcrumbPage>
            </BreadcrumbItem>
          </>
        )}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
