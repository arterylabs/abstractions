"use client";

import { useState, useEffect, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { UserRole } from "@/types/roles";
import { fetchRoles } from "@/store/modules/auth";
import { RootState, AppDispatch } from "@/store";
import { useParams } from "next/navigation";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { getGrantableRoles } from "@/lib/role-hierarchy";
import { Notifier } from "@/components/ui/sonner";

interface SendInvitationsDialogProps {
  onSendInvitations: (
    emails: string[],
    role: string,
  ) => Promise<SendInvitationResult | Error>;
  isLoading?: boolean;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

interface SendInvitationResult {
  errors?: { [email: string]: string };
  successCount?: number;
  successfulEmails?: string[];
}

type WorkspaceRoleId = "owner" | "admin" | "member";

const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

const parseInvitationEmails = (value: string): string[] =>
  value
    .split(/[\s,]+/)
    .map((email) => email.trim())
    .filter((email) => email.length > 0);

const getGrantableRoleIds = (currentUserRole?: string): WorkspaceRoleId[] => {
  if (!currentUserRole) return [];
  return getGrantableRoles(currentUserRole as WorkspaceRoleId);
};

const hasInvitationFailures = (result: SendInvitationResult): boolean =>
  Boolean(result.errors && Object.keys(result.errors).length > 0);

const formatInvitationCount = (count: number): string =>
  `${count} invitation${count > 1 ? "s" : ""}`;

export function SendInvitationsDialog({
  onSendInvitations,
  isLoading,
  open,
  onOpenChange,
}: SendInvitationsDialogProps) {
  const [emails, setEmails] = useState("");
  const [role, setRole] = useState<string>(UserRole.MEMBER);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [failedEmails, setFailedEmails] = useState<{ [email: string]: string }>(
    {},
  );

  const params = useParams();
  const dispatch = useDispatch<AppDispatch>();
  const { roles, rolesStatus, currentUser } = useSelector(
    (state: RootState) => state.auth,
  );

  const currentUserRole = currentUser?.workspaceRole;
  const availableRoles = useMemo(() => {
    const grantableRoles = getGrantableRoleIds(currentUserRole);
    return (
      roles?.filter((role) =>
        grantableRoles.includes(role.id as WorkspaceRoleId),
      ) || []
    );
  }, [currentUserRole, roles]);

  useEffect(() => {
    if (params.workspaceId && roles === null && rolesStatus === "idle") {
      dispatch(fetchRoles(params.workspaceId as string));
    }
  }, [dispatch, params.workspaceId, roles, rolesStatus]);

  useEffect(() => {
    if (availableRoles.length > 0) {
      const defaultRole =
        availableRoles.find((role) => role.id === UserRole.MEMBER) ||
        availableRoles[0];
      if (defaultRole && role !== defaultRole.id) {
        setRole(defaultRole.id);
      }
    }
  }, [availableRoles, role]);

  const handleSend = async () => {
    setError(null);
    setSuccessMessage(null);
    setFailedEmails({});

    const emailList = parseInvitationEmails(emails);

    if (emailList.length === 0) {
      setError("Please enter at least one email address");
      return;
    }

    const invalidEmails = emailList.filter((email) => !isValidEmail(email));
    if (invalidEmails.length > 0) {
      setError(`Invalid email format: ${invalidEmails.join(", ")}`);
      return;
    }

    try {
      const result = await onSendInvitations(emailList, role);
      if (result instanceof Error) {
        setError(result.message);
        return;
      }

      if (result.errors && Object.keys(result.errors).length > 0) {
        setFailedEmails(result.errors);
      }

      if (result.successCount && result.successCount > 0) {
        if (!hasInvitationFailures(result)) {
          Notifier.success({
            title: "Invitations sent successfully",
            message: `${formatInvitationCount(result.successCount)} sent successfully`,
          });

          setEmails("");
          setError(null);
          setSuccessMessage(null);
          setFailedEmails({});
          onOpenChange(false);
        } else {
          setSuccessMessage(
            `Successfully sent ${formatInvitationCount(result.successCount)}`,
          );

          if (result.successfulEmails) {
            const remainingEmails = emailList
              .filter((email) => !result.successfulEmails?.includes(email))
              .join("\n");
            setEmails(remainingEmails);
          }
        }
      }

      if (!result.successCount) {
        setError("Failed to send invitations. Please try again.");
      }
    } catch (err: Error | unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred. Please try again.");
      }
    }
  };

  const handleClose = (open: boolean) => {
    if (!open) {
      setError(null);
      setSuccessMessage(null);
      setFailedEmails({});
    }
    onOpenChange(open);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="w-full max-w-lg">
        <DialogHeader>
          <DialogTitle>Send Workspace Invitations</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 py-4">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          {successMessage && (
            <Alert>
              <AlertDescription>{successMessage}</AlertDescription>
            </Alert>
          )}
          {Object.keys(failedEmails).length > 0 && (
            <Alert variant="destructive">
              <AlertDescription>
                <div>Failed to send invitations to:</div>
                <ul className="list-disc pl-4 mt-2">
                  {Object.entries(failedEmails).map(([email, error]) => (
                    <li key={email}>
                      {email}: {error}
                    </li>
                  ))}
                </ul>
              </AlertDescription>
            </Alert>
          )}
          <div className="space-y-2">
            <label className="text-sm font-medium">Email Addresses</label>
            <Textarea
              placeholder="Enter email addresses (separated by comma, space, or new line)"
              value={emails}
              onChange={(e) => setEmails(e.target.value)}
              className="min-h-[120px]"
              data-testid="invite-emails-textarea"
            />
            <p className="text-sm text-muted-foreground">
              You can enter multiple email addresses separated by commas,
              spaces, or new lines
            </p>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Role</label>
            <Select value={role} onValueChange={setRole}>
              <SelectTrigger>
                <SelectValue>
                  {roles?.find((r) => r.id === role)?.name}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {availableRoles.map((roleItem) => (
                  <SelectItem key={roleItem.id} value={roleItem.id}>
                    <div className="flex flex-col items-start gap-1">
                      <div className="font-medium">{roleItem.name}</div>
                      <div className="text-xs text-muted-foreground whitespace-normal leading-tight">
                        {roleItem.description}
                      </div>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex justify-end pt-4">
            <Button
              onClick={handleSend}
              disabled={!emails.trim() || isLoading}
              loading={isLoading}
              data-testid="send-invitations-dialog-button"
            >
              Send Invitations
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
