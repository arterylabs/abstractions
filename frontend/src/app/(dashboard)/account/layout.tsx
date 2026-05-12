"use client";

import { usePathname } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/ui/page-header";
import Link from "next/link";
import Icon from "@/components/ui/icon";

const tabs = [
  {
    name: "General",
    href: "/account",
    value: "general",
    icon: "person",
  },
  {
    name: "Security",
    href: "/account/security",
    value: "security",
    icon: "lock",
  },
  {
    name: "Advanced",
    href: "/account/advanced",
    value: "advanced",
    icon: "tune",
  },
];

export default function UserSettings({
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
        title="Account Settings"
        description="Manage your account settings, security, and preferences."
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
