import Content, { toc } from "@/content/help/index.mdx";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";

export default function DocsPage() {
  return (
    <DocsPageWrapper slug="" toc={toc}>
      <Content />
    </DocsPageWrapper>
  );
}
