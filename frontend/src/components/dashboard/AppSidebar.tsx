"use client";

import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar";
import { NavMain } from "./NavMain";
import { NavUser } from "./NavUser";
import { WorkspaceSwitcher } from "./WorkspaceSwitcher";
import { useAppSelector } from "@/store/hooks";
import Icon from "../ui/icon";
import Logo from "../ui/logo";
import Link from "next/link";

const navItems = [
  {
    title: "Repositories",
    url: "/repositories",
    icon: <Icon type="source" className="!text-lg" />,
    testId: "nav-dashboard-link",
  },
  {
    title: "Integrations",
    url: "/integrations",
    icon: <Icon type="extension" className="!text-lg" />,
    testId: "nav-integrations-link",
  },
  {
    title: "AI",
    url: "/ai",
    icon: <Icon type="smart_toy" className="!text-lg" />,
    testId: "nav-ai-link",
  },
  {
    title: "Settings",
    url: "/settings/workspace",
    icon: <Icon type="settings" className="!text-lg" />,
    testId: "nav-settings-link",
  },
  {
    title: "Help",
    url: "/help",
    icon: <Icon type="help" className="!text-lg" />,
    testId: "nav-help-link",
  },
];

export function AppSidebar({
  disableCollapse = false,
  ...props
}: React.ComponentProps<typeof Sidebar> & {
  disableCollapse?: boolean;
}) {
  const currentUser = useAppSelector((state) => state.auth.currentUser);

  return (
    <Sidebar
      id="app-sidebar"
      className="branded-background border"
      collapsible={disableCollapse ? "none" : "icon"}
      {...props}
    >
      <SidebarHeader>
        <div className="flex items-center justify-start p-2">
          <Link href="/repositories">
            <Logo className="h-5" />
          </Link>
        </div>
        <WorkspaceSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navItems} testId="dashboard-nav" />
      </SidebarContent>
      <SidebarFooter>
        {currentUser && <NavUser user={currentUser} />}
      </SidebarFooter>
    </Sidebar>
  );
}
