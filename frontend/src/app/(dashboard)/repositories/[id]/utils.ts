const getRepositoryTabs = (repoId: string) => [
  { href: `/repositories/${repoId}/overview`, value: "overview" },
  { href: `/repositories/${repoId}/newsletters`, value: "newsletters" },
  { href: `/repositories/${repoId}/topics`, value: "topics" },
  { href: `/repositories/${repoId}/settings`, value: "settings" },
];

export function getRepositorySectionValue(pathname: string, repoId: string) {
  return (
    getRepositoryTabs(repoId).find(
      (tab) => pathname === tab.href || pathname.startsWith(`${tab.href}/`),
    )?.value ?? "overview"
  );
}
