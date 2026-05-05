"use client";

import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import UserDetailsTab from "./UserDetailsTab";
import UserWorkspacesTab from "./UserWorkspacesTab";

interface UserDetailsModalProps {
  userId: string | null;
  isOpen: boolean;
  onClose: () => void;
  onUserUpdated: () => void;
}

export default function UserDetailsModal({
  userId,
  isOpen,
  onClose,
  onUserUpdated,
}: UserDetailsModalProps) {
  const [activeTab, setActiveTab] = useState("details");

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        className="w-full sm:w-[600px] lg:w-[800px] xl:w-[900px] max-w-[90vw] overflow-y-auto p-0"
        side="right"
      >
        <SheetHeader className="px-6 py-4 border-b bg-gray-50/50">
          <SheetTitle className="text-xl">User Details</SheetTitle>
        </SheetHeader>

        <div className="px-6 py-4">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-2 mb-6">
              <TabsTrigger value="details">Details</TabsTrigger>
              <TabsTrigger value="workspaces">Workspaces</TabsTrigger>
            </TabsList>

            <TabsContent value="details" className="space-y-8 mt-0">
              <UserDetailsTab userId={userId} onUserUpdated={onUserUpdated} />
            </TabsContent>

            <TabsContent value="workspaces" className="space-y-6 mt-0">
              <UserWorkspacesTab userId={userId} />
            </TabsContent>
          </Tabs>
        </div>
      </SheetContent>
    </Sheet>
  );
}
