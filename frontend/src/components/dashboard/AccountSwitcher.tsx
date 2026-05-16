"use client";

import { UserAvatar } from "@/components/ui/user-avatar";
import { Button } from "@/components/ui/button";
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
import { useLogout } from "@/hooks/use-logout";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setCurrentWorkspace, switchWorkspace } from "@/store/modules/auth";
import { Workspace } from "@/types/workspace";
import Link from "next/link";
import { useRouter } from "next/navigation";

export function AccountSwitcher() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const logout = useLogout();

  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const workspaces = useAppSelector((state) => state.auth.workspaces);

  if (!currentUser) return null;

  const handleSwitchWorkspace = async (workspace: Workspace) => {
    if (workspace.id === currentWorkspace?.id) return;
    await dispatch(switchWorkspace(workspace.id));
    dispatch(setCurrentWorkspace(workspace));
    router.refresh();
  };

  const showWorkspaceSwitcher = workspaces && workspaces.length >= 1;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="gap-1.5"
          data-testid="account-switcher-trigger"
        >
          {currentWorkspace ? (
            <>
              <UserAvatar
                name={currentWorkspace.name}
                imageUrl={currentWorkspace.logo}
                className="h-6 w-6"
              />
              <span className="max-w-28 truncate">{currentWorkspace.name}</span>
              <span className="px-1 text-muted-foreground">/</span>
              <UserAvatar
                name={currentUser.name}
                imageUrl={currentUser.avatar}
                className="h-6 w-6"
              />
              <span className="max-w-28 truncate">{currentUser.name}</span>
            </>
          ) : (
            <>
              <UserAvatar
                name={currentUser.name}
                imageUrl={currentUser.avatar}
                className="h-6 w-6"
              />
              <span className="max-w-48 truncate">{currentUser.name}</span>
            </>
          )}
          <Icon type="unfold_more" className="text-sm text-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-64" align="end" sideOffset={4}>
        {showWorkspaceSwitcher && (
          <>
            <DropdownMenuLabel className="text-xs text-muted-foreground font-normal">
              Workspaces
            </DropdownMenuLabel>
            <DropdownMenuGroup>
              {workspaces.map((workspace) => {
                const isActive = workspace.id === currentWorkspace?.id;
                return (
                  <DropdownMenuItem
                    key={workspace.id}
                    onClick={() => handleSwitchWorkspace(workspace)}
                    className="gap-2"
                    data-testid={`workspace-item-${workspace.id}`}
                  >
                    <Icon
                      type="check"
                      className={`text-sm ${isActive ? "opacity-100" : "opacity-0"}`}
                    />
                    <span className="truncate">{workspace.name}</span>
                  </DropdownMenuItem>
                );
              })}
            </DropdownMenuGroup>
          </>
        )}
        {currentWorkspace && (
          <DropdownMenuItem asChild>
            <Link
              href="/settings/workspace"
              data-testid="workspace-settings-link"
            >
              <Icon type="settings" className="text-foreground" />
              Workspace settings
            </Link>
          </DropdownMenuItem>
        )}
        <DropdownMenuItem asChild>
          <Link href="/workspaces/create" data-testid="create-workspace-link">
            <Icon type="add" className="text-foreground" />
            Create a workspace
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuLabel className="text-xs text-muted-foreground font-normal">
          Your account
        </DropdownMenuLabel>
        <DropdownMenuGroup>
          <DropdownMenuItem asChild>
            <Link href="/account" data-testid="account-link">
              <Icon type="account_circle" className="text-foreground" />
              Account
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <Link href="/account/security" data-testid="account-security-link">
              <Icon type="lock" className="text-foreground" />
              Security
            </Link>
          </DropdownMenuItem>
          {currentUser.godMode && (
            <DropdownMenuItem asChild>
              <Link href="/god-mode">
                <Icon type="admin_panel_settings" className="text-foreground" />
                God Mode
              </Link>
            </DropdownMenuItem>
          )}
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={logout}>
          <Icon type="logout" className="text-foreground" />
          Log out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
