"use client";

import { TabbedDetailsSheet } from "@/components/ui/tabbed-details-sheet";
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
  return (
    <TabbedDetailsSheet
      isOpen={isOpen}
      onClose={onClose}
      title="User Details"
      className="sm:w-[600px] lg:w-[800px] xl:w-[900px]"
      tabs={[
        {
          value: "details",
          label: "Details",
          contentClassName: "space-y-8",
          content: (
            <UserDetailsTab userId={userId} onUserUpdated={onUserUpdated} />
          ),
        },
        {
          value: "workspaces",
          label: "Workspaces",
          contentClassName: "space-y-6",
          content: <UserWorkspacesTab userId={userId} />,
        },
      ]}
    />
  );
}
