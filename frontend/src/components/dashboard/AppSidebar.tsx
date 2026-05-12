"use client";

import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavMain } from "./NavMain";
import { NavUser } from "./NavUser";
import { useAppSelector } from "@/store/hooks";
import Icon from "../ui/icon";
import Logo from "../ui/logo";

// This is sample data.
const data = {
  navMain: [
    {
      title: "Dashboard",
      url: "/dashboard",
      icon: <Icon type="home" className="!text-lg" />,
      testId: "nav-dashboard-link",
    },
    {
      title: "Settings",
      url: "/settings/workspace",
      icon: <Icon type="settings" className="!text-lg" />,
      testId: "nav-settings-link",
    },
  ],
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
          <Logo className="h-8" />
        </div>
        {currentUser && <NavUser user={currentUser} />}
      </SidebarHeader>
      <SidebarContent>
        <NavMain testId="dashboard-nav" items={data.navMain} />
      </SidebarContent>
      <SidebarFooter></SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
