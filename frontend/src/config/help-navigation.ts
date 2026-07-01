export type HelpNavItem =
  | { type: "heading"; title: string }
  | { type?: undefined; slug: string; title: string };

export const helpNavigation: HelpNavItem[] = [
  { type: "heading", title: "Getting Started" },
  { slug: "", title: "Introduction" },
  { slug: "signing-in", title: "Signing In" },
  { slug: "accepting-invitations", title: "Invitations" },
  { slug: "workspace-setup", title: "Workspace Setup" },
  { slug: "getting_started", title: "Getting Started" },
  { type: "heading", title: "Setup" },
  { slug: "ai", title: "Configuring AI" },
  { slug: "integrations", title: "Integrations" },
  { type: "heading", title: "Features" },
  { slug: "repositories", title: "Repositories" },
  { slug: "newsletters", title: "Newsletters" },
  { slug: "subscribers", title: "Subscribers" },
  { slug: "topics", title: "Topics" },
  { type: "heading", title: "Internals" },
  { slug: "how-it-works", title: "How It Works" },
  { type: "heading", title: "Settings" },
  { slug: "workspace", title: "Workspace" },
  { slug: "account", title: "Account" },
];
