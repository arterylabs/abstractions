"use client";

import { UserAvatar } from "@/components/ui/user-avatar";
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
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { User } from "@/types/user";
import Icon from "../ui/icon";
import Link from "next/link";
import { useLogout } from "@/hooks/use-logout";

export function NavUser({ user }: { user: User }) {
  const { isMobile } = useSidebar();
  const logout = useLogout();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton size="lg" data-testid="nav-user-trigger">
              <UserAvatar
                name={user.name}
                imageUrl={user.avatar}
                className="h-8 w-8 rounded-sm"
                fallbackClassName="rounded-sm"
              />
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.name}</span>
                <span className="truncate text-xs">{user.email}</span>
              </div>
              <Icon type="unfold_more" className="ml-auto text-xl" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-sm"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <UserAvatar
                  name={user.name}
                  imageUrl={user.avatar}
                  className="h-8 w-8 rounded-sm"
                  fallbackClassName="rounded-sm"
                />
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user.name}</span>
                  <span className="truncate text-xs">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem asChild>
                <Link href="/account" data-testid="account-link">
                  <Icon type="account_circle" className="text-foreground" />
                  Account
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link
                  href="/account/security"
                  data-testid="account-security-link"
                >
                  <Icon type="lock" className="text-foreground" />
                  Security
                </Link>
              </DropdownMenuItem>
              {user.godMode && (
                <DropdownMenuItem asChild>
                  <Link href="/god-mode">
                    <Icon
                      type="admin_panel_settings"
                      className="text-violet-light"
                    />
                    <span className="flex-1">God Mode</span>
                    <span className="ml-auto text-[10px] font-medium text-violet-light opacity-80">
                      Admin
                    </span>
                  </Link>
                </DropdownMenuItem>
              )}
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={logout} data-testid="logout-menu-item">
              <Icon type="logout" className="text-foreground" />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
