"use client";

import LandingPageNavbar from "@/components/landing/Navbar";
import Footer from "@/components/landing/Footer";
import { helpNavigation } from "@/config/help-navigation";
import Icon from "@/components/ui/icon";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Link from "next/link";
import { usePathname } from "next/navigation";

function DocsSidebar({
  basePath,
  pathname,
}: {
  basePath: string;
  pathname: string;
}) {
  return (
    <Tabs value={pathname} orientation="vertical" className="w-full">
      <TabsList className="h-auto flex-col items-stretch gap-0.5 border-b-0">
        {helpNavigation.map((item, idx) =>
          item.type === "heading" ? (
            <p
              key={idx}
              className={`press-label ${idx === 0 ? "mt-0" : "mt-6"} mb-2 pl-3 text-muted-foreground/70`}
            >
              {item.title}
            </p>
          ) : (
            <TabsTrigger
              key={idx}
              value={item.slug ? `${basePath}/${item.slug}` : basePath}
              variant="sidebar"
              asChild
            >
              <Link href={item.slug ? `${basePath}/${item.slug}` : basePath}>
                {item.title}
              </Link>
            </TabsTrigger>
          ),
        )}
      </TabsList>
    </Tabs>
  );
}

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const basePath = "/docs";

  return (
    <div className="flex flex-col min-h-screen">
      <LandingPageNavbar />

      {/* Mobile top bar */}
      <div className="lg:hidden flex items-center gap-3 px-4 py-3 border-b bg-background/95 sticky top-14 z-40">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="sm" className="gap-2 text-sm">
              <Icon type="menu" className="text-base" />
              Menu
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-72 p-6 overflow-y-auto">
            <DocsSidebar basePath={basePath} pathname={pathname} />
          </SheetContent>
        </Sheet>
      </div>

      <div className="flex-1 w-full max-w-[1200px] mx-auto px-8">
        <div className="flex gap-8">
          {/* Left sidebar — hidden on mobile */}
          <aside className="hidden lg:block w-56 shrink-0">
            <div className="sticky top-[57px] h-[calc(100vh-57px)] overflow-y-auto py-8 pr-4">
              <DocsSidebar basePath={basePath} pathname={pathname} />
            </div>
          </aside>

          {/* Main content */}
          <main className="flex-1 min-w-0 py-8 lg:py-10 docs-content">
            {children}
          </main>
        </div>
      </div>
      <Footer />
    </div>
  );
}
