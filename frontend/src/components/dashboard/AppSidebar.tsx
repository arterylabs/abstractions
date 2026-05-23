"use client";

import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
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
    title: "Getting Started",
    url: "/getting-started",
    icon: <Icon type="rocket_launch" className="!text-lg" />,
    testId: "nav-getting-started-link",
  },
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
];

const helpNavItem = {
  title: "Help",
  url: "/docs",
  icon: <Icon type="help" className="!text-lg" />,
  testId: "nav-help-link",
};

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
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem className="flex justify-end">
              <SidebarMenuButton asChild variant="outline" className="">
                <Link
                  href={helpNavItem.url}
                  data-testid={helpNavItem.testId}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {helpNavItem.icon}
                  <span>{helpNavItem.title}</span>
                  <span className="ml-auto text-xs opacity-50 group-data-[collapsible=icon]:hidden">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      className="h-3 w-3"
                      aria-hidden="true"
                    >
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                      <polyline points="15 3 21 3 21 9" />
                      <line x1="10" y1="14" x2="21" y2="3" />
                    </svg>
                  </span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
        {currentUser && <NavUser user={currentUser} />}
      </SidebarFooter>
    </Sidebar>
  );
}
