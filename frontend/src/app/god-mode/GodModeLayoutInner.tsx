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
        <SidebarProvider>
          <GodModeSidebar />
          <SidebarInset className="flex flex-col overflow-hidden">
            <AdminModeBanner />
            <div className="flex w-full flex-1 flex-col gap-6 overflow-y-auto p-4 md:p-8 scrollbar-gutter-stable">
              {children}
            </div>
          </SidebarInset>
        </SidebarProvider>
      </GodModeWrapper>
    </Provider>
  );
}
