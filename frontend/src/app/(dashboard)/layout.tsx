// import Sidebar from "@/components/Sidebar";
"use client";
import Authenticated from "@/components/authenticated";
import WorkspaceContext from "@/components/WorkspaceContext";
import { AppSidebar } from "@/components/dashboard/AppSidebar";
import ImpersonationBanner from "@/components/dashboard/ImpersonationBanner";
import Icon from "@/components/ui/icon";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { store } from "@/store";
import { Provider } from "react-redux";

const DISABLE_COLLAPSE = true;

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
      <Authenticated className="h-screen w-screen overflow-hidden">
        <WorkspaceContext>
          <ImpersonationBanner />
          <SidebarProvider>
            <AppSidebar disableCollapse={DISABLE_COLLAPSE} />
            <SidebarInset>
              <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
                <div className="flex items-center gap-2 px-4">
                  {DISABLE_COLLAPSE ? null : (
                    <SidebarTrigger className="-ml-1">
                      <Icon type="menu" className="w-6 h-6" />
                    </SidebarTrigger>
                  )}
                  <Separator orientation="vertical" className="mr-2 h-4" />
                </div>
              </header>
              <div className="flex flex-1 flex-col gap-4 p-4 pt-0 max-w-3xl w-full mx-auto scrollbar-gutter-stable">
                {children}
              </div>
            </SidebarInset>
          </SidebarProvider>
        </WorkspaceContext>
      </Authenticated>
    </Provider>
  );
}
