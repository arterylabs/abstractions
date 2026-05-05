"use client";

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import UserDetailsTab from "./UserDetailsTab";

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
          <UserDetailsTab userId={userId} onUserUpdated={onUserUpdated} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
