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
      <Authenticated className="h-screen w-screen overflow-hidden">
        <WorkspaceContext>
          <SidebarProvider>
            <AppSidebar />
            <SidebarInset className="flex flex-col overflow-hidden">
              <ImpersonationBanner />
              {/* Mobile-only header so users can open the off-canvas sidebar */}
              <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4 md:hidden">
                <SidebarTrigger className="-ml-1" />
              </header>
              <div className="flex flex-1 flex-col gap-4 p-4 md:p-8 scrollbar-gutter-stable overflow-y-auto w-full max-w-[1200px] mx-auto">
                {children}
              </div>
            </SidebarInset>
          </SidebarProvider>
        </WorkspaceContext>
      </Authenticated>
    </Provider>
  );
}
