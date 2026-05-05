"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAppSelector } from "@/store/hooks";
import Icon from "../ui/icon";
import { HeaderUser } from "./HeaderUser";
import { Button } from "../ui/button";
import Logo from "../ui/logo";

const data = {
  navMain: [
    {
      title: "Home",
      url: "/dashboard",
      icon: <Icon type="home" className="!text-lg" />,
    },
    {
      title: "Subscription",
      url: "/subscription",
      icon: <Icon type="money" className="!text-lg" />,
    },
    {
      title: "Download",
      url: "/download",
      icon: <Icon type="download" className="!text-lg" />,
    },
    {
      title: "Guide",
      url: "/guide",
      icon: <Icon type="help" className="!text-lg" />,
    },
  ],
};

export function AppHeader() {
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const pathname = usePathname();

  return (
    <header
      id="#app-header"
      className="w-full min-h-16 border-b bg-[url('/assets/img/background.png')] bg-cover"
    >
      <div className="container max-w-3xl h-full mx-auto pb-2 pt-4">
        <div className="flex items-center justify-between mb-4 px-4">
          <Logo className="h-8" />
        </div>
        <div className="flex items-center justify-between">
          <div className="hidden md:flex items-center space-x-4 app-header-nav">
            {data.navMain.map((item) => (
              <Button
                variant="ghost"
                key={item.url}
                className={pathname === item.url ? "bg-accent" : ""}
              >
                <Link
                  href={item.url}
                  className="flex items-center space-x-2 text-sm"
                >
                  {item.icon}
                  <span>{item.title}</span>
                </Link>
              </Button>
            ))}
          </div>
          <div className="flex items-center">
            {currentUser && <HeaderUser user={currentUser} />}
          </div>
        </div>
      </div>
    </header>
  );
}
