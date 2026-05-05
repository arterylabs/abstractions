"use client";

import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { GodModeNavMain } from "./GodModeNavMain";
import Icon from "../ui/icon";
import Logo from "../ui/logo";
import { Button } from "@/components/ui/button";
import { NavUser } from "@/components/dashboard/NavUser";
import { useAppSelector } from "@/store/hooks";
import { useRouter } from "next/navigation";
import { useSidebar } from "@/components/ui/sidebar";

// God mode navigation data
const data = {
  navMain: [
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
  ],
};

export function GodModeSidebar({
  ...props
}: React.ComponentProps<typeof Sidebar>) {
  const router = useRouter();
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const { state } = useSidebar();

  const exitGodMode = () => {
    router.push("/dashboard");
  };

  return (
    <Sidebar
      id="god-mode-app-sidebar"
      className=""
      collapsible="icon"
      {...props}
    >
      <SidebarHeader>
        <div
          className={`flex items-center justify-start ${
            state === "expanded" ? "p-2" : "p-0"
          } py-2`}
        >
          <Logo className="h-8" icon={state !== "expanded"} />
        </div>
        <div className="px-4 py-2">
          <div
            className={`flex items-center ${
              state === "expanded" ? "justify-between" : "justify-center"
            }`}
          >
            {state === "expanded" && (
              <div>
                <h2 className="text-sm font-semibold">God Mode</h2>
                <p className="text-xs">Admin Dashboard</p>
              </div>
            )}
            <Button
              onClick={exitGodMode}
              variant="destructive"
              size="sm"
              className="h-7 px-2 text-xs"
            >
              <Icon
                type="exit_to_app"
                className={`h-3 w-3 ${state === "expanded" ? "mr-1" : ""}`}
              />
              {state === "expanded" && "Exit"}
            </Button>
          </div>
        </div>
      </SidebarHeader>
      <SidebarContent>
        <GodModeNavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        {currentUser && <NavUser user={currentUser} />}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
