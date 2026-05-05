"use client";

import { GodModeSidebar } from "@/components/god-mode/GodModeSidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { store } from "@/store";
import { Provider } from "react-redux";
import GodModeWrapper from "@/components/god-mode-wrapper";

export default function GodModeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
      <GodModeWrapper>
        <SidebarProvider>
          <GodModeSidebar />
          <SidebarInset>
            <div className="flex flex-1 flex-col gap-4 lg:py-12 py-6 lg:px-12 px-6 w-full">
              {children}
            </div>
          </SidebarInset>
        </SidebarProvider>
      </GodModeWrapper>
    </Provider>
  );
}
