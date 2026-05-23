import type React from "react";

export const helpPages: Record<
  string,
  () => Promise<{ default: React.ComponentType }>
> = {
  "getting-started": () => import("@/content/help/getting-started.mdx"),
  ai: () => import("@/content/help/ai.mdx"),
  integrations: () => import("@/content/help/integrations.mdx"),
  repositories: () => import("@/content/help/repositories.mdx"),
  newsletters: () => import("@/content/help/newsletters.mdx"),
  topics: () => import("@/content/help/topics.mdx"),
  workspace: () => import("@/content/help/workspace.mdx"),
  account: () => import("@/content/help/account.mdx"),
};
