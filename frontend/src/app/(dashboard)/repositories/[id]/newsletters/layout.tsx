"use client";

import Link from "next/link";
import { useParams, usePathname } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { getNewsletterSectionValue } from "./utils";

const getNewsletterTabs = (repoId: string) => [
  {
    name: "Issues",
    href: `/repositories/${repoId}/newsletters`,
    value: "issues",
  },
  {
    name: "Topics",
    href: `/repositories/${repoId}/newsletters/topics`,
    value: "topics",
  },
  {
    name: "Schedule",
    href: `/repositories/${repoId}/newsletters/schedule`,
    value: "schedule",
  },
  {
    name: "Recipients",
    href: `/repositories/${repoId}/newsletters/recipients`,
    value: "recipients",
  },
];

export default function NewslettersLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const tabs = getNewsletterTabs(repoId);
  const currentTab = getNewsletterSectionValue(pathname, repoId);

  return (
    <div className="flex flex-col gap-8 md:grid md:grid-cols-[160px_1fr] md:items-start md:gap-10">
      <Tabs value={currentTab} orientation="vertical">
        <TabsList className="h-auto flex-col items-stretch gap-0.5 border-b-0">
          {tabs.map((tab) => (
            <TabsTrigger
              key={tab.value}
              value={tab.value}
              variant="sidebar"
              asChild
            >
              <Link href={tab.href}>{tab.name}</Link>
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="min-w-0">{children}</div>
    </div>
  );
}
