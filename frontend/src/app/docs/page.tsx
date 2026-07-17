import type { Metadata } from "next";
import Content, { toc } from "@/content/help/index.mdx";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";
import { ogImages } from "../shared-metadata";

export const metadata: Metadata = {
  title: "Documentation",
  description:
    "Learn how to connect your repositories, configure newsletters, manage your workspace, and get the most out of Abstractions.",
  alternates: { canonical: "/docs" },
  openGraph: { type: "website", url: "/docs", images: ogImages },
};

export default function DocsPage() {
  return (
    <DocsPageWrapper slug="" toc={toc}>
      <Content />
    </DocsPageWrapper>
  );
}
