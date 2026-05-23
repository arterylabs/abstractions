import type { NextConfig } from "next";
import createMDX from "@next/mdx";

const nextConfig: NextConfig = {
  pageExtensions: ["js", "jsx", "md", "mdx", "ts", "tsx"],
  output: "standalone",
  experimental: {
    // Uses the Rust-based MDX compiler, required for Turbopack compatibility.
    // JS-based rehype/remark plugins are not supported with mdxRs.
    // GFM (tables, strikethrough, etc.) is enabled via mdxType: 'gfm'.
    mdxRs: {
      mdxType: "gfm",
    },
  },
};

const withMDX = createMDX({});

export default withMDX(nextConfig);
