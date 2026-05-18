"use client";

import { useEffect, useCallback, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "@/store";
import {
  fetchInvitations,
  resendInvitation,
  revokeInvitation,
  sendInvitations,
} from "@/store/modules/invitations";
import { AppDispatch } from "@/store";
import { InvitationCard } from "@/components/invitations/InvitationCard";
import { SendInvitationsDialog } from "@/components/invitations/SendInvitationsDialog";
import { SearchBar } from "@/components/ui/search-bar";
import { InfiniteScroll } from "@/components/ui/infinite-scroll";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";

export default function InvitationsPage() {
  const dispatch = useDispatch<AppDispatch>();
  const [sendingInvitations, setSendingInvitations] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const workspace = useSelector(
    (state: RootState) => state.auth.currentWorkspace,
  );

  const { invitations, loading, pageInfo } = useSelector(
    (state: RootState) => state.invitations,
  );

  const hasMore = pageInfo?.hasNextPage ?? false;

  useEffect(() => {
    if (workspace) {
      dispatch(
        fetchInvitations({
          workspaceId: workspace.id,
          params: {},
          append: false,
        }),
      );
    }
  }, [dispatch, workspace]);

  const handleSearch = useCallback(
    (value: string) => {
      if (!workspace) return;
      setSearchQuery(value);
      dispatch(
        fetchInvitations({
          workspaceId: workspace.id,
          params: { query: value },
          append: false,
        }),
      );
    },
    [dispatch, workspace],
  );

  const handleLoadMore = useCallback(() => {
    if (!loading && hasMore && workspace) {
      const currentPage = pageInfo?.page || 1;
      dispatch(
        fetchInvitations({
          workspaceId: workspace.id,
          params: {
            page: currentPage + 1,
            query: searchQuery || undefined,
          },
          append: true,
        }),
      );
    }
  }, [dispatch, loading, hasMore, workspace, searchQuery, pageInfo]);

  const handleResend = useCallback(
    async (invitationId: string) => {
      if (!workspace) return;
      await dispatch(
        resendInvitation({ workspaceId: workspace.id, invitationId }),
      ).unwrap();
    },
    [dispatch, workspace],
  );

  const handleRevoke = useCallback(
    async (invitationId: string) => {
      if (!workspace) return;
      await dispatch(
        revokeInvitation({ workspaceId: workspace.id, invitationId }),
      ).unwrap();
    },
    [dispatch, workspace],
  );

  const handleSendInvitations = async (
    emails: string[],
    role: string,
  ): Promise<
    | {
        errors?: { [email: string]: string };
        successCount?: number;
        successfulEmails?: string[];
      }
    | Error
  > => {
    if (!workspace) return new Error("Workspace not found");
    setSendingInvitations(true);
    try {
      const result = await dispatch(
        sendInvitations({ workspaceId: workspace.id, emails, role }),
      ).unwrap();
      dispatch(
        fetchInvitations({
          workspaceId: workspace.id,
          params: {},
          append: false,
        }),
      );
      return {
        errors: result.errors,
        successCount: result.successCount,
        successfulEmails: result.successfulEmails,
      };
    } catch (error) {
      return error instanceof Error
        ? error
        : new Error("We ran into an unexpected error");
    } finally {
      setSendingInvitations(false);
    }
  };

  if (!workspace) {
    return null;
  }

  return (
    <div className="space-y-8">
      <div>
        <div className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <SearchBar
              placeholder="Search invitations..."
              onSearch={handleSearch}
              className="w-full"
            />
            <Button onClick={() => setDialogOpen(true)}>
              Send Invitations
            </Button>
          </div>

          <InfiniteScroll
            hasMore={hasMore}
            loading={loading}
            onLoadMore={handleLoadMore}
          >
            <div className="grid gap-4">
              {(invitations || []).map((invitation) => (
                <InvitationCard
                  key={invitation.id}
                  invitation={{
                    id: invitation.id,
                    email: invitation.email,
                    role: invitation.role,
                    state: invitation.state,
                    createdAt: invitation.createdAt,
                    invitedBy: {
                      name: invitation.invitedBy?.name || "Unknown",
                      email: invitation.invitedBy?.email || "",
                      avatar: invitation.invitedBy?.avatar,
                    },
                  }}
                  onResend={handleResend}
                  onRevoke={handleRevoke}
                />
              ))}

              {!loading &&
                invitations.length === 0 &&
                (searchQuery ? (
                  <EmptyState
                    icon="search_off"
                    title="No invitations found"
                    description="No invitations match your search. Try a different query."
                    cta={null}
                  />
                ) : (
                  <EmptyState
                    icon="mail"
                    title="No invitations yet"
                    description="Send invitations to add members to your workspace."
                    cta={
                      <Button onClick={() => setDialogOpen(true)}>
                        Send Invitations
                      </Button>
                    }
                  />
                ))}
            </div>
          </InfiniteScroll>

          <SendInvitationsDialog
            open={dialogOpen}
            onOpenChange={setDialogOpen}
            onSendInvitations={handleSendInvitations}
            isLoading={sendingInvitations}
          />
        </div>
      </div>
    </div>
  );
}
