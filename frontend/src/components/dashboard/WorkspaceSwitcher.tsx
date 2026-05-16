"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { getInitials } from "@/lib/user";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setCurrentWorkspace, switchWorkspace } from "@/store/modules/auth";
import { Workspace } from "@/types/workspace";
import Link from "next/link";

export function WorkspaceSwitcher() {
  const { isMobile } = useSidebar();
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const workspaces = useAppSelector((state) => state.auth.workspaces);

  if (!currentWorkspace) return null;

  const handleSwitchWorkspace = async (workspace: Workspace) => {
    if (workspace.id === currentWorkspace.id) return;
    await dispatch(switchWorkspace(workspace.id));
    dispatch(setCurrentWorkspace(workspace));
  };

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              data-testid="workspace-switcher-trigger"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage
                  src={currentWorkspace.logo ?? undefined}
                  alt={currentWorkspace.name}
                />
                <AvatarFallback
                  className="rounded-lg"
                  name={currentWorkspace.name}
                >
                  {getInitials(currentWorkspace.name)}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">
                  {currentWorkspace.name}
                </span>
              </div>
              <Icon type="unfold_more" className="ml-auto text-xl" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="start"
            sideOffset={4}
          >
            {workspaces && workspaces.length > 0 && (
              <>
                <DropdownMenuLabel className="text-xs text-muted-foreground font-normal">
                  Workspaces
                </DropdownMenuLabel>
                <DropdownMenuGroup>
                  {workspaces.map((workspace) => {
                    const isActive = workspace.id === currentWorkspace.id;
                    return (
                      <DropdownMenuItem
                        key={workspace.id}
                        onClick={() => handleSwitchWorkspace(workspace)}
                        className="gap-2"
                        data-testid={`workspace-item-${workspace.id}`}
                      >
                        <Avatar className="h-6 w-6 rounded-md">
                          <AvatarImage
                            src={workspace.logo ?? undefined}
                            alt={workspace.name}
                          />
                          <AvatarFallback
                            className="rounded-md text-xs"
                            name={workspace.name}
                          >
                            {getInitials(workspace.name)}
                          </AvatarFallback>
                        </Avatar>
                        <span className="truncate">{workspace.name}</span>
                        {isActive && (
                          <Icon type="check" className="ml-auto text-sm" />
                        )}
                      </DropdownMenuItem>
                    );
                  })}
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
              </>
            )}
            <DropdownMenuGroup>
              <DropdownMenuItem asChild>
                <Link
                  href="/settings/workspace"
                  data-testid="workspace-settings-link"
                >
                  <Icon type="settings" className="text-foreground" />
                  Workspace settings
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link
                  href="/workspaces/create"
                  data-testid="create-workspace-link"
                >
                  <Icon type="add" className="text-foreground" />
                  Create workspace
                </Link>
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
