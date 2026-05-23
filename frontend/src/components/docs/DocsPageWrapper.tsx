import { type ReactNode } from "react";
import DocsPager from "@/components/docs/DocsPager";
import DocsBreadcrumbs from "@/components/docs/DocsBreadcrumbs";
import TableOfContents, {
  type TocEntry,
} from "@/components/docs/TableOfContents";

interface DocsPageWrapperProps {
  children: ReactNode;
  toc?: TocEntry[];
  slug: string;
  basePath?: "/docs" | "/help";
  showBreadcrumbs?: boolean;
}

export default function DocsPageWrapper({
  children,
  toc = [],
  slug,
  basePath = "/docs",
  showBreadcrumbs = true,
}: DocsPageWrapperProps) {
  return (
    <div className="flex gap-10 w-full">
      {/* Main content + pager */}
      <div className="flex-1 min-w-0">
        {showBreadcrumbs && <DocsBreadcrumbs slug={slug} basePath={basePath} />}
        {children}
        <DocsPager slug={slug} basePath={basePath} />
      </div>

      {/* On-page TOC — only shown on XL screens when there are headings */}
      {toc.length > 0 && (
        <aside className="hidden xl:block w-44 shrink-0">
          <div className="sticky top-[57px] pt-10 max-h-[calc(100vh-57px)] overflow-y-auto">
            <TableOfContents toc={toc} />
          </div>
        </aside>
      )}
    </div>
  );
}
