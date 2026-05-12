"use client";

import { useEffect, useState } from "react";
import { isAxiosError } from "axios";
import api from "@/api";
import useAPI from "@/api/hooks";
import { NewsletterCard } from "@/components/newsletters/NewsletterCard";
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
import { getErrorMessage } from "@/lib/errors";
import { useParams, useRouter } from "next/navigation";

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
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [isGeneratingDraft, setIsGeneratingDraft] = useState(false);
  const [isTriggeringNewsletter, setIsTriggeringNewsletter] = useState(false);

  useEffect(() => {
    if (!repoId) {
      return;
    }

    void loadNewsletters(repoId);
    void loadTopics(repoId);
  }, [loadNewsletters, loadTopics, repoId]);

  const handleGenerateDraft = async () => {
    setIsGeneratingDraft(true);

    try {
      await api.newsletters.generatePreview(repoId);
      await Promise.all([loadNewsletters(repoId), loadTopics(repoId)]);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to generate draft",
        message: getErrorMessage(error),
      });
    } finally {
      setIsGeneratingDraft(false);
    }
  };

  const handleTriggerNewsletter = async () => {
    setIsTriggeringNewsletter(true);

    try {
      await api.newsletters.triggerNewsletter(repoId);
      await Promise.all([loadNewsletters(repoId), loadTopics(repoId)]);
      Notifier.success({
        title: "Newsletter queued",
        message: "Newsletter generation has been queued successfully.",
      });
    } catch (error: unknown) {
      const isRateLimited =
        isAxiosError(error) && error.response?.status === 429;
      Notifier.error({
        title: "Unable to trigger newsletter",
        message: isRateLimited
          ? "Newsletter already sent in the last 24 hours."
          : getErrorMessage(error),
      });
    } finally {
      setIsTriggeringNewsletter(false);
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Newsletter History"
        showBackButton
        handleBack={() => router.push("/repositories")}
      />

      <Card>
        <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <CardTitle>Next Newsletter</CardTitle>
            <CardDescription>
              {topicsResponse.result?.length
                ? `${topicsResponse.result.length} topic${topicsResponse.result.length === 1 ? "" : "s"} queued for the next draft.`
                : "No pending topics queued yet."}
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button
              variant="outline"
              onClick={() => void handleTriggerNewsletter()}
              loading={isTriggeringNewsletter}
            >
              Send Newsletter Now
            </Button>
            <Button
              onClick={() => void handleGenerateDraft()}
              loading={isGeneratingDraft}
            >
              Generate Draft
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {topicsResponse.isFetching && !topicsResponse.result ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Spinner className="h-5 w-5" />
              Loading queued topics...
            </div>
          ) : topicsResponse.result?.length ? (
            <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
              {topicsResponse.result.map((topic) => (
                <li key={topic.id}>
                  <span className="font-medium text-foreground">
                    {topicToLabel(topic)}
                  </span>{" "}
                  <span>{topic.sentAt ? "Sent" : "Pending"}</span>
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
            <NewsletterCard
              key={newsletter.id}
              newsletter={newsletter}
              onClick={() =>
                router.push(
                  `/repositories/${repoId}/newsletters/${newsletter.id}`,
                )
              }
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title="No newsletters yet"
          description="Connected repositories will show newsletter history here once your first send goes out."
          cta={
            <Button
              variant="outline"
              onClick={() => router.push("/repositories")}
            >
              Back to Repositories
            </Button>
          }
        />
      )}
    </div>
  );
}
