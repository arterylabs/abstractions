import { StatusBadge } from "@/components/ui/status-badge";
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
import { getErrorMessage } from "@/lib/errors";

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
            message: getErrorMessage(error),
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
            message: getErrorMessage(error),
          });
        }
      },
    });
  };

  const getRoleName = (roleId: string) => {
    const role = roles?.find((r) => r.id === roleId);
    return role?.name || roleId;
  };

  return (
    <Card className="press-card w-full overflow-hidden">
      <CardContent className="p-5">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-grow">
            <div className="flex flex-wrap items-center gap-3">
              <span className="font-mono text-sm text-foreground">
                {invitation.email}
              </span>
              <StatusBadge status={invitation.state} />
            </div>
            <span className="press-label mt-2 flex items-center gap-1.5 text-muted-foreground">
              <Icon type="assignment_ind" className="!text-sm" />
              {getRoleName(invitation.role)}
            </span>
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
        <p className="mt-4 border-t border-border pt-3 font-mono text-xs text-muted-foreground">
          Invited by{" "}
          <span className="text-foreground">{invitation.invitedBy.name}</span> ·{" "}
          {formatDistanceToNow(new Date(invitation.createdAt))} ago
        </p>
      </CardContent>
    </Card>
  );
}
