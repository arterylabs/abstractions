"use client";

import { usePathname } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import Link from "next/link";

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
    icon: "mail",
  },
  {
    name: "Billing",
    href: "/settings/workspace/billing",
    value: "billing",
    icon: "payment",
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
    <div className="">
      <PageHeader
        title="Workspace Settings"
        description="Manage your workspace settings, users, and billing information."
      >
        <Tabs value={currentTab}>
          <TabsList>
            {tabs.map((tab) => (
              <TabsTrigger key={tab.value} value={tab.value} asChild>
                <Link href={tab.href}>
                  <Icon type={tab.icon} className="text-lg" />
                  {tab.name}
                </Link>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </PageHeader>
      <div className="pt-6">{children}</div>
    </div>
  );
}
