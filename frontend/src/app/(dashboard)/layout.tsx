"use client";
import Authenticated from "@/components/authenticated";
import WorkspaceContext from "@/components/WorkspaceContext";
import { DashboardTopNav } from "@/components/dashboard/DashboardTopNav";
import ImpersonationBanner from "@/components/dashboard/ImpersonationBanner";
import { store } from "@/store";
import { Provider } from "react-redux";

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
          <div className="flex flex-col h-full">
            <DashboardTopNav />
            <div className="flex flex-1 flex-col gap-4 p-4 max-w-6xl w-full mx-auto scrollbar-gutter-stable overflow-y-auto">
              {children}
            </div>
          </div>
        </WorkspaceContext>
      </Authenticated>
    </Provider>
  );
}
