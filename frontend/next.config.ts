import type { NextConfig } from "next";

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

export default nextConfig;
