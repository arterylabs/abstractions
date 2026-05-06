"use client";

import { useEffect, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { formatDateTime } from "@/lib/format";
import { getErrorMessage } from "@/lib/errors";
import { useParams, useRouter } from "next/navigation";

export function formatRate(rate: number): string {
  return `${(rate * 100).toFixed(1)}%`;
}

const topicToLabel = (topic: unknown) => {
  if (typeof topic === "string") {
    return topic;
  }

  if (topic && typeof topic === "object") {
    const topicRecord = topic as Record<string, unknown>;
    if (typeof topicRecord.title === "string") {
      return topicRecord.title;
    }
    if (typeof topicRecord.name === "string") {
      return topicRecord.name;
    }
    return JSON.stringify(topicRecord);
  }

  return "Untitled topic";
};

export default function RepositoryNewslettersPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;

  const [newslettersResponse, loadNewsletters] = useAPI(api.newsletters.list);
  const [upcomingResponse, loadUpcoming] = useAPI(api.newsletters.getUpcoming);
  const [isRegenerating, setIsRegenerating] = useState(false);

  useEffect(() => {
    if (!repoId) {
      return;
    }

    void loadNewsletters(repoId);
    void loadUpcoming(repoId);
  }, [loadNewsletters, loadUpcoming, repoId]);

  const handleRegenerateTopics = async () => {
    setIsRegenerating(true);

    try {
      await api.newsletters.regenerateTopics(repoId);
      await loadUpcoming(repoId);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to regenerate topics",
        message: getErrorMessage(error),
      });
    } finally {
      setIsRegenerating(false);
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Newsletter History"
        showBackButton
        handleBack={() => router.push("/dashboard/repositories")}
      />

      <Card>
        <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <CardTitle>Next Newsletter</CardTitle>
            <CardDescription>
              {upcomingResponse.result?.nextSendAt
                ? `Scheduled for ${formatDateTime(upcomingResponse.result.nextSendAt)}`
                : "No upcoming newsletter scheduled yet."}
            </CardDescription>
          </div>
          <Button
            onClick={() => void handleRegenerateTopics()}
            loading={isRegenerating}
          >
            Regenerate Topics
          </Button>
        </CardHeader>
        <CardContent>
          {upcomingResponse.isFetching && !upcomingResponse.result ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Spinner className="h-5 w-5" />
              Loading next newsletter...
            </div>
          ) : upcomingResponse.result?.topics?.length ? (
            <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
              {upcomingResponse.result.topics.map((topic, index) => (
                <li key={`${topicToLabel(topic)}-${index}`}>
                  {topicToLabel(topic)}
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-muted-foreground">
              No topics generated yet.
            </p>
          )}
        </CardContent>
      </Card>

      {newslettersResponse.isFetching && !newslettersResponse.result ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading newsletters...
          </p>
        </div>
      ) : newslettersResponse.result &&
        newslettersResponse.result.length > 0 ? (
        <div className="grid gap-4">
          {newslettersResponse.result.map((newsletter) => (
            <button
              key={newsletter.id}
              type="button"
              className="text-left"
              onClick={() =>
                router.push(
                  `/dashboard/repositories/${repoId}/newsletters/${newsletter.id}`,
                )
              }
            >
              <Card className="transition-colors hover:border-primary/40">
                <CardHeader>
                  <CardTitle className="text-lg">
                    {formatDateTime(
                      newsletter.sentAt ?? newsletter.createdAt,
                    ) ?? "Draft newsletter"}
                  </CardTitle>
                </CardHeader>
                <CardContent className="grid gap-3 text-sm text-muted-foreground sm:grid-cols-3">
                  <div>
                    <span className="font-medium text-foreground">
                      Recipients:
                    </span>{" "}
                    {newsletter.recipientCount}
                  </div>
                  <div>
                    <span className="font-medium text-foreground">
                      Open rate:
                    </span>{" "}
                    {formatRate(newsletter.openRate)}
                  </div>
                  <div>
                    <span className="font-medium text-foreground">
                      Click rate:
                    </span>{" "}
                    {formatRate(newsletter.clickRate)}
                  </div>
                </CardContent>
              </Card>
            </button>
          ))}
        </div>
      ) : (
        <EmptyState
          title="No newsletters yet"
          description="Connected repositories will show newsletter history here once your first send goes out."
          cta={
            <Button
              variant="outline"
              onClick={() => router.push("/dashboard/repositories")}
            >
              Back to Repositories
            </Button>
          }
        />
      )}
    </div>
  );
}
