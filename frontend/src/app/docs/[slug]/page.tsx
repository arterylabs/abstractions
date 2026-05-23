import { notFound } from "next/navigation";
import { helpPages } from "@/config/help-pages";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";
import type { TocEntry } from "@/components/ui/table-of-contents";

export default async function DocsSlugPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const loader = helpPages[slug];
  if (!loader) notFound();

  const mod = await loader();
  const Content = mod.default;
  const toc: TocEntry[] = (mod as { toc?: TocEntry[] }).toc ?? [];

  return (
    <DocsPageWrapper slug={slug} toc={toc}>
      <Content />
    </DocsPageWrapper>
  );
}
