"use client";

import { useEffect } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { useNewsletterPolling } from "@/hooks/use-newsletter-polling";
import { NewsletterCard } from "@/components/newsletters/NewsletterCard";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import { SectionMark } from "@/components/ui/section-mark";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { describeRecipients } from "@/lib/repository";
import { useParams, useRouter } from "next/navigation";
import { useAppSelector } from "@/store/hooks";

/* The archive. Sending is a masthead action and the queue lives in Topics, so
   this page is only concerned with issues that already exist. */
export default function RepositoryNewslettersPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [newslettersResponse, loadNewsletters] = useNewsletterPolling(
    currentWorkspace?.id,
    repoId,
  );
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);

  useEffect(() => {
    if (!repoId || !currentWorkspace?.id) {
      return;
    }

    void loadTopics(currentWorkspace.id, repoId);
  }, [loadTopics, repoId, currentWorkspace?.id]);

  const newsletters = newslettersResponse.result ?? [];
  const queuedCount = (topicsResponse.result ?? []).filter(
    (topic) => !topic.sentAt,
  ).length;

  const handleResend = async (newsletterId: string) => {
    try {
      await api.newsletters.resendNewsletter(
        currentWorkspace!.id,
        repoId,
        newsletterId,
      );
      await loadNewsletters(currentWorkspace!.id, repoId);
      Notifier.info({
        title: "Newsletter re-queued",
        message: "The newsletter has been queued for resending.",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to resend newsletter",
        message: getErrorMessage(error),
      });
    }
  };

  return (
    <div className="flex flex-col gap-9">
      <PageHeader
        size="h2"
        rule={false}
        title="Newsletters"
        description="Every issue that has gone out for this repository. Open one to read it as your team received it, or re-queue it to send again."
      />

      {/* A single line locating this page between the queue and the
          recipients — both one click away in the sidebar next door. */}
      <p className="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm text-muted-foreground">
        <span>
          {queuedCount > 0
            ? `${queuedCount} ${queuedCount === 1 ? "topic" : "topics"} queued for the next issue`
            : "Nothing queued for the next issue"}
        </span>
        {repository && (
          <>
            <span aria-hidden className="text-muted-foreground/40">
              ·
            </span>
            <span>Going out to {describeRecipients(repository)}</span>
          </>
        )}
      </p>

      {newslettersResponse.isFetching && newsletters.length === 0 ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="press-label text-muted-foreground">
            Pulling the archive…
          </p>
        </div>
      ) : newsletters.length > 0 ? (
        <div className="flex flex-col gap-4">
          <SectionMark
            label="Back issues"
            action={
              <span className="press-label text-muted-foreground">
                {newsletters.length}{" "}
                {newsletters.length === 1 ? "issue" : "issues"}
              </span>
            }
          />
          {newsletters.map((newsletter) => (
            <NewsletterCard
              key={newsletter.id}
              newsletter={newsletter}
              onClick={() =>
                router.push(
                  `/repositories/${repoId}/newsletters/${newsletter.id}`,
                )
              }
              onResend={() => handleResend(newsletter.id)}
            />
          ))}
        </div>
      ) : (
        <EmptyState
          icon="newspaper"
          title="Your send archive is empty"
          description="Sent newsletters are filed here. Queue a topic and send your first issue from the masthead above."
          cta={
            <Button
              variant="outline"
              onClick={() =>
                router.push(`/repositories/${repoId}/newsletters/topics`)
              }
            >
              <Icon type="subject" className="!text-base" />
              Go to topics
            </Button>
          }
        />
      )}
    </div>
  );
}
