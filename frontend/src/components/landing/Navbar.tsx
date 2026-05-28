"use client";

import Link from "next/link";
import Logo from "@/components/ui/logo";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import Icon from "@/components/ui/icon";
import api from "@/api";
import { useEffect, useState } from "react";
// import Image from "next/image"; // used by GitHub button — restore when open source launch

const NavLinks = ({ mobile = false }: { mobile?: boolean }) => {
  const links = [
    { href: "/#features", label: "Features" },
    { href: "/#how-it-works", label: "How it works" },
    { href: "/#faq", label: "FAQs" },
    { href: "/docs", label: "Docs" },
  ];

  return (
    <div
      className={`${
        mobile ? "flex flex-col space-y-4" : "hidden lg:flex items-center gap-6"
      }`}
    >
      {links.map((link) => (
        <Link key={link.href} href={link.href}>
          <Button variant="ghost">{link.label}</Button>
        </Link>
      ))}
    </div>
  );
};

export default function Navbar() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    setIsLoggedIn(!!api.utils.getToken());
  }, []);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container h-14 max-w-[1200px] mx-auto flex items-center justify-between px-8">
        <Link className="flex items-center space-x-2" href="/#">
          <Logo className="h-8 w-auto" />
        </Link>

        <NavLinks />

        <div className="flex items-center gap-4">
          {isLoggedIn ? (
            <>
              <Link href="/dashboard" className="hidden lg:block">
                <Button>My Account</Button>
              </Link>
            </>
          ) : (
            <>
              {/* GitHub navbar link — hidden until open source launch */}
              {/* <Link
                href="https://github.com/nshntarora/abstractions"
                className="hidden lg:block"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button variant="ghost">
                  <Image
                    src="/assets/img/github.svg"
                    alt="GitHub"
                    width={16}
                    height={16}
                    className="mr-1 invert"
                  />
                  GitHub
                </Button>
              </Link> */}
              <Link href="/login" className="hidden lg:block">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link href="/signup" className="hidden lg:block">
                <Button>Sign up</Button>
              </Link>
            </>
          )}

          <Sheet>
            <SheetTrigger asChild className="lg:hidden">
              <Button variant="ghost" size="icon">
                <Icon type="menu" className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent>
              <div className="flex flex-col space-y-4 mt-8">
                <NavLinks mobile />
                {isLoggedIn ? (
                  <Link href="/dashboard">
                    <Button className="w-full">My Account</Button>
                  </Link>
                ) : (
                  <>
                    {/* GitHub mobile nav link — hidden until open source launch */}
                    {/* <Link
                      href="https://github.com/nshntarora/abstractions"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <Button variant="ghost" className="w-full">
                        <Icon type="code" className="text-base mr-1" />
                        GitHub
                      </Button>
                    </Link> */}
                    <Link href="/login">
                      <Button variant="ghost" className="w-full">
                        Login
                      </Button>
                    </Link>
                    <Link href="/signup">
                      <Button className="w-full">Sign up</Button>
                    </Link>
                  </>
                )}
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  );
}
