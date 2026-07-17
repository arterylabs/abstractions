import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { helpPages } from "@/config/help-pages";
import { helpNavigation } from "@/config/help-navigation";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";
import type { TocEntry } from "@/components/ui/table-of-contents";
import { ogImages, siteUrl } from "../../shared-metadata";

type Props = { params: Promise<{ slug: string }> };

function getPageTitle(slug: string): string {
  const navItem = helpNavigation.find(
    (item) => item.type !== "heading" && item.slug === slug,
  );
  return navItem && navItem.type !== "heading" ? navItem.title : slug;
}

export async function generateStaticParams() {
  return Object.keys(helpPages).map((slug) => ({ slug }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  if (!helpPages[slug]) return {};

  const pageTitle = getPageTitle(slug);

  return {
    title: pageTitle,
    description: `Abstractions documentation: ${pageTitle}. Learn how to use this feature in your engineering workflow.`,
    alternates: { canonical: `/docs/${slug}` },
    openGraph: { type: "article", url: `/docs/${slug}`, images: ogImages },
  };
}

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
  const pageTitle = getPageTitle(slug);

  const techArticleSchema = {
    "@context": "https://schema.org",
    "@type": "TechArticle",
    headline: pageTitle,
    description: `Abstractions documentation: ${pageTitle}.`,
    url: `${siteUrl}/docs/${slug}`,
    isPartOf: { "@type": "WebSite", url: siteUrl, name: "Abstractions" },
  };

  return (
    <DocsPageWrapper slug={slug} toc={toc}>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(techArticleSchema) }}
      />
      <Content />
    </DocsPageWrapper>
  );
}
