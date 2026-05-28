import type { Metadata } from "next";
import Content, { toc } from "@/content/help/index.mdx";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";

export const metadata: Metadata = {
  title: "Documentation",
  description:
    "Learn how to connect your repositories, configure newsletters, manage your workspace, and get the most out of Abstractions.",
  alternates: { canonical: "/docs" },
  openGraph: { type: "website", url: "/docs" },
};

export default function DocsPage() {
  return (
    <DocsPageWrapper slug="" toc={toc}>
      <Content />
    </DocsPageWrapper>
  );
}
