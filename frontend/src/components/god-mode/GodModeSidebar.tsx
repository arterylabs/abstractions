"use client";

import * as React from "react";
import Link from "next/link";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavMain } from "@/components/dashboard/NavMain";
import { NavUser } from "@/components/dashboard/NavUser";
import Icon from "@/components/ui/icon";
import Logo from "@/components/ui/logo";
import { Badge } from "@/components/ui/badge";
import { useAppSelector } from "@/store/hooks";
import { useSidebar } from "@/components/ui/sidebar";
import { cn } from "@/lib/utils";

const navItems = [
  {
    title: "Dashboard",
    url: "/god-mode",
    icon: <Icon type="dashboard" className="!text-lg" />,
  },
  {
    title: "Users",
    url: "/god-mode/users",
    icon: <Icon type="group" className="!text-lg" />,
  },
  {
    title: "Workspaces",
    url: "/god-mode/workspaces",
    icon: <Icon type="workspaces" className="!text-lg" />,
  },
  {
    title: "Configuration",
    url: "/god-mode/configuration",
    icon: <Icon type="tune" className="!text-lg" />,
  },
];

export function GodModeSidebar({
  ...props
}: React.ComponentProps<typeof Sidebar>) {
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const { state } = useSidebar();
  const expanded = state === "expanded";

  return (
    <Sidebar
      id="god-mode-app-sidebar"
      className={cn(
        "border-r border-sidebar-border",
        // Violet admin accent: left edge + active/hover nav states
        "border-l-2 border-l-violet",
        "[&_[data-slot=sidebar-menu-button][data-active=true]]:bg-violet/15",
        "[&_[data-slot=sidebar-menu-button][data-active=true]]:text-violet-light",
        "[&_[data-slot=sidebar-menu-button][data-active=true]_span.text-primary]:text-violet-light",
        "[&_[data-slot=sidebar-menu-button]:hover]:bg-violet/10",
      )}
      collapsible="icon"
      {...props}
    >
      <SidebarHeader>
        <div className="flex items-center justify-start gap-2 p-2">
          <Link href="/god-mode" className="flex items-center gap-2 min-w-0">
            <Logo className="h-5" height="auto" width="auto" />
          </Link>
          {expanded && (
            <Badge
              variant="outline"
              size="sm"
              className="border-violet/40 bg-violet/10 text-violet-light shrink-0"
            >
              Admin
            </Badge>
          )}
        </div>
        {!expanded && (
          <div className="flex justify-center pb-1" title="Admin mode">
            <Icon
              type="admin_panel_settings"
              className="text-base text-violet-light"
            />
          </div>
        )}
      </SidebarHeader>
      <SidebarContent>
        <NavMain
          items={navItems}
          exactMatch
          testId="god-mode-nav"
          activeIconClassName="text-violet-light"
        />
      </SidebarContent>
      <SidebarFooter>
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild variant="outline">
                <Link href="/repositories">
                  <Icon type="arrow_back" className="!text-lg" />
                  <span>Back to app</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
        {currentUser && <NavUser user={currentUser} />}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
