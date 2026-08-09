"use client";

import { usePathname } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/ui/page-header";
import Link from "next/link";
import Icon from "@/components/ui/icon";

const tabs = [
  {
    name: "General",
    href: "/settings/workspace",
    value: "general",
    icon: "settings",
  },
  {
    name: "Users",
    href: "/settings/workspace/users",
    value: "users",
    icon: "group",
  },
  {
    name: "Invitations",
    href: "/settings/workspace/invitations",
    value: "invitations",
    icon: "person_add",
  },
];

export default function WorkspaceSettings({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  const currentTab =
    tabs.find((tab) => tab.href === pathname)?.value || "general";

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        eyebrow="Workspace"
        rule={false}
        title="Workspace settings"
        description="Configure your workspace, manage your team, and set up AI models."
      />
      <div className="flex flex-col gap-6">
        <Tabs value={currentTab}>
          <TabsList>
            {tabs.map((tab) => (
              <TabsTrigger key={tab.value} value={tab.value} asChild>
                <Link href={tab.href}>
                  <Icon type={tab.icon} className="!text-sm" />
                  {tab.name}
                </Link>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
        {children}
      </div>
    </div>
  );
}
