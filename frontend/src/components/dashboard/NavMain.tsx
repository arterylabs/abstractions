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
}: {
  title?: string;
  testId?: string;
  items: {
    title: string;
    url: string;
    icon?: React.ReactNode;
    testId?: string;
  }[];
}) {
  const pathname = usePathname();

  return (
    <SidebarGroup data-testid={testId}>
      {title && <SidebarGroupLabel>{title}</SidebarGroupLabel>}
      <SidebarMenu>
        {items.map((item) => {
          const isActive = pathname.startsWith(item.url);

          return (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild isActive={isActive} variant="outline">
                <Link href={item.url} data-testid={item.testId}>
                  {item.icon && (
                    <span className={isActive ? "text-primary" : undefined}>
                      {item.icon}
                    </span>
                  )}
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}
