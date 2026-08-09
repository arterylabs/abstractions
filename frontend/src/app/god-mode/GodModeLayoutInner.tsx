"use client";

import { GodModeSidebar } from "@/components/god-mode/GodModeSidebar";
import { AdminModeBanner } from "@/components/god-mode/AdminModeBanner";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { store } from "@/store";
import { Provider } from "react-redux";
import GodModeWrapper from "@/components/god-mode-wrapper";

export default function GodModeLayoutInner({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
      <GodModeWrapper>
        {/* Same viewport-locked shell as the dashboard: document never scrolls;
            the main pane owns overflow with a stable scrollbar gutter. */}
        <SidebarProvider className="h-svh overflow-hidden">
          <GodModeSidebar />
          <SidebarInset className="flex min-h-0 flex-col overflow-hidden">
            <AdminModeBanner />
            <div className="scrollbar-gutter-stable min-h-0 flex-1 overflow-y-auto">
              <div className="flex w-full flex-col gap-6 p-4 md:p-8">
                {children}
              </div>
            </div>
          </SidebarInset>
        </SidebarProvider>
      </GodModeWrapper>
    </Provider>
  );
}
