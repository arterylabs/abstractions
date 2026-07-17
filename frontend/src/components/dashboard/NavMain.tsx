"use client";

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import Link from "next/link";
import { usePathname } from "next/navigation";

export function NavMain({
  title,
  items,
  testId,
  exactMatch = false,
  activeIconClassName = "text-primary",
}: {
  title?: string;
  testId?: string;
  exactMatch?: boolean;
  /** Class applied to the icon when the item is active */
  activeIconClassName?: string;
  items: {
    title: string;
    url: string;
    icon?: React.ReactNode;
    testId?: string;
    external?: boolean;
  }[];
}) {
  const pathname = usePathname();

  return (
    <SidebarGroup data-testid={testId}>
      {title && <SidebarGroupLabel>{title}</SidebarGroupLabel>}
      <SidebarMenu>
        {items.map((item) => {
          const isActive =
            !item.external &&
            (exactMatch
              ? pathname === item.url
              : pathname.startsWith(item.url));

          return (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild isActive={isActive} variant="outline">
                <Link
                  href={item.url}
                  data-testid={item.testId}
                  {...(item.external
                    ? { target: "_blank", rel: "noopener noreferrer" }
                    : {})}
                >
                  {item.icon && (
                    <span
                      className={isActive ? activeIconClassName : undefined}
                    >
                      {item.icon}
                    </span>
                  )}
                  <span>{item.title}</span>
                  {item.external && (
                    <span className="ml-auto opacity-50 text-xs">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        className="w-3 h-3"
                        aria-hidden="true"
                      >
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                        <polyline points="15 3 21 3 21 9" />
                        <line x1="10" y1="14" x2="21" y2="3" />
                      </svg>
                    </span>
                  )}
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}
