"use client";

import { usePathname } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import Link from "next/link";

const tabs = [
  {
    name: "General",
    href: "/account",
    value: "general",
    icon: "account_circle",
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
    <div className="">
      <PageHeader
        title="Account Settings"
        description="Manage your account settings, security, and preferences."
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
