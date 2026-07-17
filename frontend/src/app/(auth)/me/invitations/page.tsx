"use client";

import { useCallback, useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { UserAvatar } from "@/components/ui/user-avatar";
import { useRouter } from "next/navigation";
import { formatDistanceToNow } from "date-fns";
import { Invitation } from "@/api/invitations";
import { Notifier } from "@/components/ui/sonner";
import api from "@/api";
import Spinner from "@/components/ui/spinner";
import AuthWrapper from "@/components/authenticated";
import { InfiniteScroll } from "@/components/ui/infinite-scroll";
import { PageInfo } from "@/types/pagination";
import { EmptyState } from "@/components/ui/empty-state";

function InvitationsPageContent() {
  const [invitations, setInvitations] = useState<Invitation[]>([]);
  const [loadingInvitations, setLoadingInvitations] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [pageInfo, setPageInfo] = useState<PageInfo | null>(null);
  const [processingInvitationId, setProcessingInvitationId] = useState<
    string | null
  >(null);
  const router = useRouter();

  const fetchInvitations = useCallback(async (page = 1, append = false) => {
    try {
      if (page === 1) {
        setLoadingInvitations(true);
      } else {
        setLoadingMore(true);
      }

      const response = await api.invitations.getUserPendingInvitations({
        page,
        perPage: 10,
      });

      const newInvitations = response.data.results;

      if (append) {
        setInvitations((prev) => [...prev, ...newInvitations]);
      } else {
        setInvitations(newInvitations);
      }

      setPageInfo(response.data.pageInfo);
    } catch (err) {
      Notifier.error({
        title: "Error",
        message: (err as Error).message || "Failed to load invitations",
      });
    } finally {
      setLoadingInvitations(false);
      setLoadingMore(false);
    }
  }, []);

  const handleLoadMore = useCallback(() => {
    if (pageInfo && pageInfo.hasNextPage && !loadingMore) {
      const nextPage = (pageInfo.page || 1) + 1;
      fetchInvitations(nextPage, true);
    }
  }, [pageInfo, loadingMore, fetchInvitations]);

  useEffect(() => {
    fetchInvitations();
  }, [fetchInvitations]);

  const handleAcceptInvitation = async (invitation: Invitation) => {
    setProcessingInvitationId(invitation.id);
    try {
      await api.invitations.acceptInvitation(invitation.id);
      setInvitations((prev) => prev.filter((inv) => inv.id !== invitation.id));
      Notifier.success({
        title: "Invitation Accepted",
        message: "You have been added to the workspace",
      });
      router.push("/dashboard");
    } catch (error) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message || "Failed to accept invitation",
      });
    } finally {
      setProcessingInvitationId(null);
    }
  };

  const handleIgnoreInvitation = async (invitation: Invitation) => {
    setProcessingInvitationId(invitation.id);
    try {
      await api.invitations.ignoreInvitation(invitation.id);
      setInvitations((prev) => prev.filter((inv) => inv.id !== invitation.id));
      Notifier.success({
        title: "Invitation Ignored",
        message: "The invitation has been removed",
      });
    } catch (error) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message || "Failed to ignore invitation",
      });
    } finally {
      setProcessingInvitationId(null);
    }
  };

  if (loadingInvitations) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <Spinner />
      </div>
    );
  }

  if (invitations.length === 0) {
    return (
      <EmptyState
        icon="mail"
        title="No Pending Invitations"
        description="You don't have any pending invitations at the moment."
        cta={null}
      />
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-semibold mb-6">Pending Invitations</h1>
      <InfiniteScroll
        hasMore={pageInfo?.hasNextPage ?? false}
        loading={loadingMore}
        onLoadMore={handleLoadMore}
      >
        <div className="grid gap-4">
          {invitations.map((invitation) => (
            <Card key={invitation.id}>
              <CardContent className="p-6">
                <div className="flex flex-col sm:flex-row items-start gap-4">
                  <UserAvatar
                    name={invitation.invitedBy?.name || ""}
                    imageUrl={invitation.invitedBy?.avatar}
                    className="h-12 w-12 shrink-0"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-lg">
                      {invitation.invitedBy?.name || "Unknown"} invited you to
                      join {invitation.workspaceName}
                    </div>
                    <div className="text-base mb-1 text-muted-foreground">
                      You&apos;ve been invited as a{" "}
                      {invitation.role === "admin"
                        ? "Workspace Admin"
                        : "Workspace Member"}
                    </div>
                    <div className="text-sm text-muted-foreground/80 mb-4">
                      Invitation sent by {invitation.invitedBy?.email} • Expires{" "}
                      {formatDistanceToNow(new Date(invitation.expiresAt), {
                        addSuffix: true,
                      })}
                    </div>
                    <div className="flex flex-wrap gap-3">
                      <Button
                        onClick={() => handleAcceptInvitation(invitation)}
                        disabled={processingInvitationId === invitation.id}
                      >
                        Accept Invitation
                      </Button>
                      <Button
                        variant="outline"
                        onClick={() => handleIgnoreInvitation(invitation)}
                        disabled={processingInvitationId === invitation.id}
                      >
                        Ignore
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </InfiniteScroll>
    </div>
  );
}

export default function AcceptInvitationPage() {
  return (
    <AuthWrapper>
      <InvitationsPageContent />
    </AuthWrapper>
  );
}
