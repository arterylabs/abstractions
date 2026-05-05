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
      title: "Events",
      url: "/events",
      icon: <Icon type="event" className="!text-lg" />,
    },
    {
      title: "Help",
      url: "/help",
      icon: <Icon type="help" className="!text-lg" />,
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
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter></SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
