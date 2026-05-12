"use client";

import { useEffect } from "react";
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
import { formatDateTime } from "@/lib/format";
import { useParams, useRouter } from "next/navigation";

export function formatRate(rate: number): string {
  return `${(rate * 100).toFixed(1)}%`;
}

export default function RepositoryNewslettersPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;

  const [newslettersResponse, loadNewsletters] = useAPI(api.newsletters.list);
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);

  useEffect(() => {
    if (!repoId) {
      return;
    }

    void loadNewsletters(repoId);
    void loadTopics(repoId);
  }, [loadNewsletters, loadTopics, repoId]);

  return (
    <div className="space-y-6">
      <PageHeader
        title="Newsletter History"
        showBackButton
        handleBack={() => router.push("/dashboard/repositories")}
      />

      <Card>
        <CardHeader>
          <CardTitle>Newsletter Topics</CardTitle>
          <CardDescription>
            Pending topics stay in the queue until a newsletter is sent.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {topicsResponse.isFetching && !topicsResponse.result ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Spinner className="h-5 w-5" />
              Loading topics...
            </div>
          ) : topicsResponse.result?.length ? (
            <ul className="space-y-3 text-sm">
              {topicsResponse.result.map((topic) => (
                <li key={topic.id} className="rounded-md border p-3">
                  <div className="flex items-center justify-between gap-3">
                    <span className="font-medium text-foreground">
                      {topic.title}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {topic.sentAt ? "Sent" : "Pending"}
                    </span>
                  </div>
                  {topic.description ? (
                    <p className="mt-2 text-muted-foreground">
                      {topic.description}
                    </p>
                  ) : null}
                  <p className="mt-2 text-xs text-muted-foreground">
                    {topic.sentAt
                      ? `Sent ${formatDateTime(topic.sentAt)}`
                      : `Queued ${formatDateTime(topic.createdAt)}`}
                  </p>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-muted-foreground">
              No topics queued yet.
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
