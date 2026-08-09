"use client";

import Authenticated from "@/components/authenticated";
import WorkspaceContext from "@/components/WorkspaceContext";
import { AppSidebar } from "@/components/dashboard/AppSidebar";
import ImpersonationBanner from "@/components/dashboard/ImpersonationBanner";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { store } from "@/store";
import { Provider } from "react-redux";

export default function DashboardLayoutInner({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
      {/* Loading shell only — Authenticated applies className to the spinner
          state, not the authenticated tree. Height is locked on the provider. */}
      <Authenticated className="h-svh w-full overflow-hidden">
        <WorkspaceContext>
          {/* Lock the shell to the viewport so the document never scrolls.
              Scroll lives on the main pane with scrollbar-gutter: stable so
              the gutter is always reserved and content doesn't shift when
              pages of different heights gain/lose a scrollbar. */}
          <SidebarProvider className="h-svh overflow-hidden">
            <AppSidebar />
            <SidebarInset className="flex min-h-0 flex-col overflow-hidden">
              <ImpersonationBanner />
              {/* Mobile-only header so users can open the off-canvas sidebar */}
              <header className="flex h-12 shrink-0 items-center gap-3 border-b border-border px-4 md:hidden">
                <SidebarTrigger className="-ml-1" />
                <span className="press-label text-muted-foreground">
                  Abstractions
                </span>
              </header>
              {/* Full-width scroller: gutter sits at the pane edge, not inside
                  the max-width column. min-h-0 lets the flex child shrink so
                  overflow-y actually engages. */}
              <div className="scrollbar-gutter-stable min-h-0 flex-1 overflow-y-auto">
                <div className="mx-auto flex w-full max-w-[1180px] flex-col gap-8 p-4 md:p-10">
                  {children}
                </div>
              </div>
            </SidebarInset>
          </SidebarProvider>
        </WorkspaceContext>
      </Authenticated>
    </Provider>
  );
}
