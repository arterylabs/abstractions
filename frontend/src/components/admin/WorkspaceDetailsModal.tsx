"use client";

import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import WorkspaceDetailsTab from "./WorkspaceDetailsTab";
import WorkspaceUsersTab from "./WorkspaceUsersTab";
import WorkspaceSubscriptionTab from "./WorkspaceSubscriptionTab";

interface WorkspaceDetailsModalProps {
  workspaceId: string | null;
  isOpen: boolean;
  onClose: () => void;
  onWorkspaceUpdated: () => void;
}

export default function WorkspaceDetailsModal({
  workspaceId,
  isOpen,
  onClose,
  onWorkspaceUpdated,
}: WorkspaceDetailsModalProps) {
  const [activeTab, setActiveTab] = useState("details");

  const handleWorkspaceDeleted = () => {
    onWorkspaceUpdated();
    onClose();
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        className="w-full sm:w-[700px] lg:w-[900px] xl:w-[1000px] max-w-[95vw] overflow-y-auto p-0"
        side="right"
      >
        <SheetHeader className="px-6 py-4 border-b bg-gray-50/50">
          <SheetTitle className="text-xl">Workspace Details</SheetTitle>
        </SheetHeader>

        <div className="px-6 py-4">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-3 mb-6">
              <TabsTrigger value="details">Details</TabsTrigger>
              <TabsTrigger value="users">Users</TabsTrigger>
              <TabsTrigger value="subscription">Subscription</TabsTrigger>
            </TabsList>

            <TabsContent value="details" className="space-y-8 mt-0">
              <WorkspaceDetailsTab
                workspaceId={workspaceId}
                onWorkspaceUpdated={onWorkspaceUpdated}
                onWorkspaceDeleted={handleWorkspaceDeleted}
              />
            </TabsContent>

            <TabsContent value="users" className="space-y-6 mt-0">
              <WorkspaceUsersTab workspaceId={workspaceId} />
            </TabsContent>

            <TabsContent value="subscription" className="space-y-6 mt-0">
              <WorkspaceSubscriptionTab workspaceId={workspaceId} />
            </TabsContent>
          </Tabs>
        </div>
      </SheetContent>
    </Sheet>
  );
}
