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
  {
    name: "Billing",
    href: "/settings/workspace/billing",
    value: "billing",
    icon: "payments",
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
    <div className="space-y-6">
      <PageHeader
        title="Workspace Settings"
        description="Manage your workspace settings, users, and billing information."
      />
      <div className="flex flex-row gap-6 items-start">
        <Tabs value={currentTab} orientation="vertical">
          <TabsList className="flex-col h-auto w-44 bg-transparent justify-start shrink-0 p-0">
            {tabs.map((tab) => (
              <TabsTrigger
                key={tab.value}
                value={tab.value}
                variant="sidebar"
                className="justify-start"
                asChild
              >
                <Link href={tab.href}>
                  <Icon type={tab.icon} className="!text-base" />
                  {tab.name}
                </Link>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
        <div className="flex-1">{children}</div>
      </div>
    </div>
  );
}
