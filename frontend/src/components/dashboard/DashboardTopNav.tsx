"use client";

import * as React from "react";
import Link from "next/link";
// import { usePathname } from "next/navigation";
import { useAppSelector } from "@/store/hooks";
// import Icon from "../ui/icon";
import Logo from "../ui/logo";
import { HeaderUser } from "./HeaderUser";
// import { Button } from "../ui/button";

// const navItems = [
//   {
//     title: "Events",
//     url: "/events",
//     icon: <Icon type="event" className="!text-lg" />,
//   },
//   {
//     title: "Help",
//     url: "/help",
//     icon: <Icon type="help" className="!text-lg" />,
//   },
// ];

export function DashboardTopNav() {
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  // const pathname = usePathname();

  return (
    <header className="w-full h-12 border-b branded-background shrink-0">
      <div className="mx-auto w-full max-w-6xl px-6 h-full flex items-center">
        <div className="flex items-center gap-4 flex-1">
          <Link href="/events">
            <Logo className="h-6 w-auto" />
          </Link>
          {/* <nav className="flex items-center gap-1" data-testid="dashboard-nav">
            {navItems.map((item) => {
              const isActive = pathname.startsWith(item.url);

              return (
                <Button
                  key={item.url}
                  variant="ghost"
                  asChild
                  className={isActive ? "bg-accent" : ""}
                  data-testid={`nav-${item.title.toLowerCase()}-link`}
                >
                  <Link href={item.url} className="flex items-center gap-2">
                    <span className={isActive ? "text-primary" : undefined}>
                      {item.icon}
                    </span>
                    <span>{item.title}</span>
                  </Link>
                </Button>
              );
            })}
          </nav> */}
        </div>
        <div className="flex items-center">
          {currentUser && <HeaderUser user={currentUser} />}
        </div>
      </div>
    </header>
  );
}
