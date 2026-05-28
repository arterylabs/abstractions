export type HelpNavItem =
  | { type: "heading"; title: string }
  | { type?: undefined; slug: string; title: string };

export const helpNavigation: HelpNavItem[] = [
  { type: "heading", title: "Getting Started" },
  { slug: "", title: "Introduction" },
  { slug: "workspace-setup", title: "Setting Up Your Workspace" },
  { slug: "getting-started", title: "Getting Started" },
  { type: "heading", title: "Setup" },
  { slug: "ai", title: "Configuring AI" },
  { slug: "integrations", title: "Integrations" },
  { type: "heading", title: "Features" },
  { slug: "repositories", title: "Repositories" },
  { slug: "newsletters", title: "Newsletters" },
  { slug: "subscribers", title: "Managing Subscribers" },
  { slug: "topics", title: "Topics" },
  { type: "heading", title: "Internals" },
  { slug: "how-it-works", title: "How It Works" },
  { type: "heading", title: "Settings" },
  { slug: "workspace", title: "Workspace Settings" },
  { slug: "account", title: "Account & Security" },
];
