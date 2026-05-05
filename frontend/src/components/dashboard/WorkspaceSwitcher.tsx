import * as React from "react";
import { useRouter } from "next/navigation";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  useSidebar,
} from "@/components/ui/sidebar";
import Icon from "../ui/icon";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setCurrentWorkspace } from "@/store/modules/auth";
import { Workspace } from "@/types/workspace";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import Link from "next/link";

export function WorkspaceSwitcher() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const { isMobile } = useSidebar();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const workspaces = useAppSelector((state) => state.auth.workspaces);

  const handleSwitchWorkspace = async (workspace: Workspace) => {
    dispatch(setCurrentWorkspace(workspace));
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n.charAt(0))
      .join("");
  };

  React.useEffect(() => {
    // Redirect if there are no workspaces
    if (workspaces && workspaces.length === 0) {
      router.push("/workspaces/create");
    }
  }, [workspaces, router]);

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              variant="primary"
              // className="border shadow-sm hover:bg-background border-border bg-sidebar-accent text-sidebar-accent-foreground hover:text-sidebar-accent-foreground"
            >
              {currentWorkspace ? (
                <div className="flex items-center gap-2 w-full justify-between">
                  <Avatar className="h-8 w-8 rounded-lg">
                    <AvatarImage
                      src={currentWorkspace.logo || undefined}
                      alt={currentWorkspace.name}
                    />
                    <AvatarFallback className="rounded-lg">
                      {getInitials(currentWorkspace.name)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="line-clamp-1 flex-1 text-sm font-medium">
                    {currentWorkspace.name}
                  </span>
                  <Icon type="unfold_more" className="ml-auto text-2xl" />
                </div>
              ) : (
                <div className="flex items-center justify-between w-full">
                  <span className="line-clamp-1 flex-1 text-sm font-medium">
                    Select Workspace
                  </span>
                  <Icon type="unfold_more" className="ml-auto text-2xl" />
                </div>
              )}
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel>
              <div className="font-normal text-sm">Your Workspaces</div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup className="max-h-[300px] overflow-y-auto">
              {workspaces?.map((workspace) => (
                <DropdownMenuItem
                  key={workspace.id}
                  onSelect={() => handleSwitchWorkspace(workspace)}
                >
                  <Avatar className="h-8 w-8 rounded-lg mr-2">
                    <AvatarImage
                      src={workspace.logo || undefined}
                      alt={workspace.name}
                    />
                    <AvatarFallback className="rounded-lg">
                      {getInitials(workspace.name)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="line-clamp-1 flex-1">{workspace.name}</span>
                  {currentWorkspace?.id === workspace.id && (
                    <Icon type="check" className="ml-2 text-lg" />
                  )}
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link href="/workspaces/create">
                <Icon type="add" className="text-foreground text-lg" />
                Create Workspace
              </Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
