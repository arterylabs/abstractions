export type HelpNavItem =
  | { type: "heading"; title: string }
  | { type?: undefined; slug: string; title: string };

export const helpNavigation: HelpNavItem[] = [
  { type: "heading", title: "Getting Started" },
  { slug: "", title: "Introduction" },
  { slug: "getting-started", title: "Getting Started" },
];
