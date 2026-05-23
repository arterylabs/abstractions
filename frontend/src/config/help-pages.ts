import type React from "react";

export const helpPages: Record<
  string,
  () => Promise<{ default: React.ComponentType }>
> = {
  "getting-started": () => import("@/content/help/getting-started.mdx"),
};
