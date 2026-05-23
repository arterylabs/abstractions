import Content, { toc } from "@/content/help/index.mdx";
import DocsPageWrapper from "@/components/docs/DocsPageWrapper";

export default function HelpPage() {
  return (
    <DocsPageWrapper slug="" toc={toc} basePath="/help">
      <Content />
    </DocsPageWrapper>
  );
}
