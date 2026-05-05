import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { askConfirmation } from "@/components/ui/confirmation";
import Icon from "@/components/ui/icon";
import { Notifier } from "@/components/ui/sonner";
import { formatDistanceToNow } from "date-fns";
import { useAppSelector } from "@/store/hooks";

interface InvitationCardProps {
  invitation: {
    id: string;
    email: string;
    role: string;
    state: string;
    invitedBy: {
      name: string;
      email: string;
      avatar?: string;
    };
    createdAt: string;
  };
  onResend: (id: string) => Promise<void>;
  onRevoke: (id: string) => Promise<void>;
}

export function InvitationCard({
  invitation,
  onResend,
  onRevoke,
}: InvitationCardProps) {
  const roles = useAppSelector((state) => state.auth.roles);

  const handleResend = () => {
    askConfirmation({
      title: "Resend Invitation",
      content: `Are you sure you want to resend the invitation to ${invitation.email}? This will send a new invitation email.`,
      variant: "primary",
      okText: "Resend Invitation",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await onResend(invitation.id);
          Notifier.success({
            title: "Invitation Resent",
            message: `Successfully resent invitation to ${invitation.email}`,
          });
        } catch (error) {
          console.error("Failed to resend invitation:", error);
          Notifier.error({
            title: "Failed to Resend Invitation",
            message:
              error instanceof Error
                ? error.message
                : "An unexpected error occurred while resending the invitation",
          });
        }
      },
    });
  };

  const handleRevoke = () => {
    askConfirmation({
      title: "Revoke Invitation",
      content: `Are you sure you want to revoke the invitation for ${invitation.email}? This action cannot be undone and the invitation link will no longer work.`,
      variant: "destructive",
      okText: "Revoke Invitation",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await onRevoke(invitation.id);
          Notifier.success({
            title: "Invitation Revoked",
            message: `Successfully revoked invitation for ${invitation.email}`,
          });
        } catch (error) {
          console.error("Failed to revoke invitation:", error);
          Notifier.error({
            title: "Failed to Revoke Invitation",
            message:
              error instanceof Error
                ? error.message
                : "An unexpected error occurred while revoking the invitation",
          });
        }
      },
    });
  };

  const stateColors = {
    pending: "bg-yellow-100 text-yellow-800",
    expired: "bg-red-100 text-red-800",
    accepted: "bg-green-100 text-green-800",
    ignored: "bg-gray-100 text-gray-800",
  };

  const getStatusDisplay = (state: string) => {
    const statusMap = {
      pending: "Pending",
      expired: "Expired",
      accepted: "Accepted",
      ignored: "Ignored",
    };
    return statusMap[state as keyof typeof statusMap] || state;
  };

  const getRoleName = (roleId: string) => {
    const role = roles?.find((r) => r.id === roleId);
    return role?.name || roleId;
  };

  return (
    <Card className="w-full">
      <CardContent className="pt-6">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-grow">
            <div className="flex items-center gap-2">
              <div className="text-md font-medium">{invitation.email}</div>
              <Badge
                variant="secondary"
                className={
                  stateColors[invitation.state as keyof typeof stateColors]
                }
              >
                {getStatusDisplay(invitation.state)}
              </Badge>
            </div>
            <div className="flex items-center text-sm text-muted-foreground mt-1">
              <Icon type="assignment_ind" className="inline mr-1" />
              {getRoleName(invitation.role)}
            </div>
          </div>
          <div className="flex items-center gap-2">
            {invitation.state === "pending" && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="min-h-8 min-w-8 p-0">
                    <Icon type="more_vert" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={handleResend}>
                    Resend
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    variant="destructive"
                    onClick={handleRevoke}
                  >
                    Revoke
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
        <div className="mt-4">
          <p className="text-sm text-muted-foreground">
            Invited by{" "}
            <span className="font-medium text-foreground">
              {invitation.invitedBy.name}
            </span>{" "}
            <span className="text-sm text-muted-foreground">
              {formatDistanceToNow(new Date(invitation.createdAt))} ago
            </span>
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
