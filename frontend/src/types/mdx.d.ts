// Augment all MDX modules to allow the `toc` named export
// (exported from individual .mdx files for on-page table of contents)
declare module "*.mdx" {
  export const toc:
    | Array<{ id: string; label: string; level: 2 | 3 }>
    | undefined;
}
