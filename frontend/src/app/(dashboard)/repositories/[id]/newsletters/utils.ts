const getNewsletterTabs = (repoId: string) => [
  { href: `/repositories/${repoId}/newsletters`, value: "issues" },
  { href: `/repositories/${repoId}/newsletters/topics`, value: "topics" },
  { href: `/repositories/${repoId}/newsletters/schedule`, value: "schedule" },
  {
    href: `/repositories/${repoId}/newsletters/recipients`,
    value: "recipients",
  },
];

/* Unlike the parent tab bar, this can't prefix-match: /newsletters is a
   real page (Issues) that every other sub-route sits underneath, so only
   an exact match picks out Topics, Schedule, or Recipients — anything
   else, including an issue's own /newsletters/<id>, falls back to Issues. */
export function getNewsletterSectionValue(pathname: string, repoId: string) {
  return (
    getNewsletterTabs(repoId).find(
      (tab) => tab.value !== "issues" && pathname === tab.href,
    )?.value ?? "issues"
  );
}
