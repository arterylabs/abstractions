import type { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
  const baseUrl =
    "https://" + (process.env.NEXT_PUBLIC_SITE_URL ?? "useabstractions.com");
  return {
    rules: [
      {
        userAgent: "*",
        allow: ["/", "/docs", "/docs/", "/privacy", "/terms", "/cookies"],
        disallow: [
          "/dashboard/",
          "/god-mode/",
          "/workspaces/",
          "/settings/",
          "/login",
          "/signup",
          "/forgot",
          "/reset",
          "/verify",
          "/app_login",
          "/authenticated",
          "/me/",
          "/api/",
          "/_next/",
        ],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
  };
}
