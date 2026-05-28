import type { MDXComponents } from "mdx/types";
import { Children, isValidElement } from "react";
import type {
  ReactNode,
  HTMLAttributes,
  ThHTMLAttributes,
  TdHTMLAttributes,
} from "react";
import Callout from "@/components/ui/callout";
import CopyButton from "@/components/ui/copy-button";

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
          {code && (
            <CopyButton
              text={code}
              className="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity bg-accent hover:bg-accent-light border border-border text-foreground/60 hover:text-foreground"
            />
          )}
        </div>
      );
    },

    // Tables — filter whitespace text nodes which are invalid direct children of <table>
    table: ({ children, ...props }: HTMLAttributes<HTMLTableElement>) => {
      const validChildren = Children.toArray(children).filter(isValidElement);
      return (
        <div className="overflow-x-auto my-4">
          <table {...props}>{validChildren}</table>
        </div>
      );
    },
    thead: ({
      children,
      ...props
    }: HTMLAttributes<HTMLTableSectionElement>) => {
      const validChildren = Children.toArray(children).filter(isValidElement);
      return <thead {...props}>{validChildren}</thead>;
    },
    tbody: ({
      children,
      ...props
    }: HTMLAttributes<HTMLTableSectionElement>) => {
      const validChildren = Children.toArray(children).filter(isValidElement);
      return <tbody {...props}>{validChildren}</tbody>;
    },
    tr: ({ children, ...props }: HTMLAttributes<HTMLTableRowElement>) => {
      const validChildren = Children.toArray(children).filter(isValidElement);
      return <tr {...props}>{validChildren}</tr>;
    },
    th: ({ children, ...props }: ThHTMLAttributes<HTMLTableCellElement>) => (
      <th {...props}>{children}</th>
    ),
    td: ({ children, ...props }: TdHTMLAttributes<HTMLTableCellElement>) => (
      <td {...props}>{children}</td>
    ),
  };
}
