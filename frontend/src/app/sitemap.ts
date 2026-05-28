import type { MetadataRoute } from "next";
import { helpPages } from "@/config/help-pages";

export default function sitemap(): MetadataRoute.Sitemap {
  const baseUrl =
    "https://" + (process.env.NEXT_PUBLIC_SITE_URL ?? "useabstractions.com");
  const now = new Date();

  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 1.0,
    },
    {
      url: `${baseUrl}/docs`,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 0.9,
    },
    {
      url: `${baseUrl}/privacy`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseUrl}/terms`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseUrl}/cookies`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
  ];

  const docPages: MetadataRoute.Sitemap = Object.keys(helpPages).map(
    (slug) => ({
      url: `${baseUrl}/docs/${slug}`,
      lastModified: now,
      changeFrequency: "weekly" as const,
      priority: 0.8,
    }),
  );

  return [...staticPages, ...docPages];
}
