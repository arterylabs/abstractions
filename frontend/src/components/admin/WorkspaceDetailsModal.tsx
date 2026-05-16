"use client";

import { TabbedDetailsSheet } from "@/components/ui/tabbed-details-sheet";
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
  const handleWorkspaceDeleted = () => {
    onWorkspaceUpdated();
    onClose();
  };

  return (
    <TabbedDetailsSheet
      isOpen={isOpen}
      onClose={onClose}
      title="Workspace Details"
      className="sm:w-[700px] lg:w-[900px] xl:w-[1000px] max-w-[95vw]"
      tabs={[
        {
          value: "details",
          label: "Details",
          contentClassName: "space-y-8",
          content: (
            <WorkspaceDetailsTab
              workspaceId={workspaceId}
              onWorkspaceUpdated={onWorkspaceUpdated}
              onWorkspaceDeleted={handleWorkspaceDeleted}
            />
          ),
        },
        {
          value: "users",
          label: "Users",
          contentClassName: "space-y-6",
          content: <WorkspaceUsersTab workspaceId={workspaceId} />,
        },
        {
          value: "subscription",
          label: "Subscription",
          contentClassName: "space-y-6",
          content: <WorkspaceSubscriptionTab workspaceId={workspaceId} />,
        },
      ]}
    />
  );
}
