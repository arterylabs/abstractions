import type React from "react";

export const helpPages: Record<
  string,
  () => Promise<{ default: React.ComponentType }>
> = {
  "signing-in": () => import("@/content/help/signing-in.mdx"),
  "accepting-invitations": () =>
    import("@/content/help/accepting-invitations.mdx"),
  "workspace-setup": () => import("@/content/help/workspace-setup.mdx"),
  getting_started: () => import("@/content/help/getting_started.mdx"),
  ai: () => import("@/content/help/ai.mdx"),
  integrations: () => import("@/content/help/integrations.mdx"),
  repositories: () => import("@/content/help/repositories.mdx"),
  newsletters: () => import("@/content/help/newsletters.mdx"),
  subscribers: () => import("@/content/help/subscribers.mdx"),
  topics: () => import("@/content/help/topics.mdx"),
  workspace: () => import("@/content/help/workspace.mdx"),
  account: () => import("@/content/help/account.mdx"),
  "how-it-works": () => import("@/content/help/how-it-works.mdx"),
};
