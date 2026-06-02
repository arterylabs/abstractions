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
    optimizePackageImports: [
      "@radix-ui/react-icons",
      "lucide-react",
      "recharts",
    ],
  },
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          {
            key: "Strict-Transport-Security",
            value: "max-age=63072000; includeSubDomains; preload",
          },
          {
            key: "Cross-Origin-Opener-Policy",
            value: "same-origin",
          },
          {
            key: "X-Content-Type-Options",
            value: "nosniff",
          },
          {
            key: "Referrer-Policy",
            value: "strict-origin-when-cross-origin",
          },
        ],
      },
    ];
  },
};

const withMDX = createMDX({});

export default withMDX(nextConfig);
