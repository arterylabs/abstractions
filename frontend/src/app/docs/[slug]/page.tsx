import { notFound } from "next/navigation";
import { helpPages } from "@/config/help-pages";

export default async function DocsSlugPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const loader = helpPages[slug];
  if (!loader) notFound();

  const { default: Content } = await loader();
  return <Content />;
}
