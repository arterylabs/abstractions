"use client";

import Authenticated from "@/components/authenticated";
import WorkspaceContext from "@/components/WorkspaceContext";
import { AppSidebar } from "@/components/dashboard/AppSidebar";
import ImpersonationBanner from "@/components/dashboard/ImpersonationBanner";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
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
