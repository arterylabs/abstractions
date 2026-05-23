import type { MDXComponents } from "mdx/types";
import type {
  ReactNode,
  HTMLAttributes,
  ThHTMLAttributes,
  TdHTMLAttributes,
} from "react";
import Callout from "@/components/docs/Callout";
import CopyButton from "@/components/docs/CopyButton";

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return {
    ...components,

    // Blockquotes become Note callouts by default
    blockquote: ({ children }: { children?: ReactNode }) => (
      <Callout type="note">{children}</Callout>
    ),

    // Code blocks with copy button
    pre: ({ children, ...props }: HTMLAttributes<HTMLPreElement>) => {
      // Extract raw text content for the copy button
      let code = "";
      const codeEl = (children as React.ReactElement<{ children?: ReactNode }>)
        ?.props?.children;
      if (typeof codeEl === "string") code = codeEl;

      return (
        <div className="relative group">
          <pre {...props}>{children}</pre>
          {code && <CopyButton text={code} />}
        </div>
      );
    },

    // Tables
    table: ({ children, ...props }: HTMLAttributes<HTMLTableElement>) => (
      <div className="overflow-x-auto my-4">
        <table {...props}>{children}</table>
      </div>
    ),
    th: ({ children, ...props }: ThHTMLAttributes<HTMLTableCellElement>) => (
      <th {...props}>{children}</th>
    ),
    td: ({ children, ...props }: TdHTMLAttributes<HTMLTableCellElement>) => (
      <td {...props}>{children}</td>
    ),
  };
}
