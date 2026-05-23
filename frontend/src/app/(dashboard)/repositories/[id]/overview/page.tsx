"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import Icon from "@/components/ui/icon";
import Spinner from "@/components/ui/spinner";
import { StatusBadge } from "@/components/ui/status-badge";
import { formatDaysAgo } from "@/lib/format";
import { useAppSelector } from "@/store/hooks";
import { NewsletterPromptCard } from "@/components/repositories/NewsletterPromptCard";

export default function RepositoryOverviewPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const repoId = params.id;

  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [newslettersResponse, loadNewsletters] = useAPI(api.newsletters.list);
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [statsResponse, loadStats] = useAPI(api.repositories.getStats);

  useEffect(() => {
    if (!repoId || !currentWorkspace?.id) return;
    void loadNewsletters(currentWorkspace.id, repoId);
    void loadTopics(currentWorkspace.id, repoId);
    void loadStats(currentWorkspace.id, repoId);
  }, [repoId, currentWorkspace?.id, loadNewsletters, loadTopics, loadStats]);

  const newsletters = newslettersResponse.result ?? [];
  const topics = topicsResponse.result ?? [];
  const stats = statsResponse.result;

  const pendingTopics = topics.filter((t) => !t.sentAt);
  const sentNewsletters = newsletters.filter((n) => n.sentAt).slice(0, 5);

  const emails = repository?.emails ?? [];
  const sendToWorkspaceMembers = repository?.sendToWorkspaceMembers ?? true;

  return (
    <div className="space-y-4">
      <NewsletterPromptCard
        repositoryId={repoId}
        newsletterEnabled={repository?.newsletterEnabled ?? false}
      />

      {/* Top row: Next Newsletter + Recipients */}
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Next Newsletter */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Next Newsletter
            </CardTitle>
            <Icon type="mail" className="!text-base text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {topicsResponse.isFetching ? (
              <div className="flex items-center gap-2 py-2">
                <Spinner className="h-4 w-4" />
                <span className="text-sm text-muted-foreground">
                  Loading topics…
                </span>
              </div>
            ) : pendingTopics.length === 0 ? (
              <EmptyState
                icon="topic"
                title="Your next newsletter is empty"
                description="Generate a draft to start building topics for the next send."
                cta={
                  <Link
                    href={`/repositories/${repoId}/newsletters`}
                    className="text-sm text-primary hover:underline"
                  >
                    Go to Newsletters
                  </Link>
                }
              />
            ) : (
              <div className="space-y-2">
                <p className="text-xs text-muted-foreground">
                  {pendingTopics.length}{" "}
                  {pendingTopics.length === 1 ? "topic" : "topics"} queued
                </p>
                <ul className="space-y-1">
                  {pendingTopics.map((topic) => (
                    <li
                      key={topic.id}
                      className="flex items-start gap-2 text-sm"
                    >
                      <Icon
                        type="fiber_manual_record"
                        className="!text-[0.5rem] mt-1.5 shrink-0 text-primary"
                      />
                      <span className="line-clamp-1">{topic.title}</span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recipients */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Recipients
            </CardTitle>
            <Icon type="group" className="!text-base text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {sendToWorkspaceMembers && (
              <Badge variant="secondary" className="mb-2">
                All workspace members
              </Badge>
            )}
            {emails.length === 0 && !sendToWorkspaceMembers ? (
              <EmptyState
                icon="group"
                title="No recipients set up"
                description="Add your team in Settings to start receiving newsletters."
                cta={
                  <Link
                    href={`/repositories/${repoId}/settings`}
                    className="text-sm text-primary hover:underline"
                  >
                    Go to Settings
                  </Link>
                }
              />
            ) : (
              <div className="flex flex-wrap gap-1.5">
                {emails.slice(0, 5).map((email) => (
                  <Badge key={email} variant="outline" size="sm">
                    {email}
                  </Badge>
                ))}
                {emails.length > 5 && (
                  <Badge variant="secondary" size="sm">
                    +{emails.length - 5} more
                  </Badge>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Indexing Stats */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Indexing
          </CardTitle>
          <Icon type="database" className="!text-base text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {statsResponse.isFetching ? (
            <div className="flex items-center gap-2 py-2">
              <Spinner className="h-4 w-4" />
              <span className="text-sm text-muted-foreground">
                Loading stats…
              </span>
            </div>
          ) : !stats || stats.fileCount === 0 ? (
            <div className="flex items-center gap-3">
              {repository && <StatusBadge status={repository.status} />}
              <span className="text-sm text-muted-foreground">
                {repository?.status === "indexing"
                  ? "Indexing in progress…"
                  : "No files indexed yet."}
              </span>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div>
                <p className="text-3xl font-semibold tabular-nums">
                  {stats.fileCount.toLocaleString()}
                </p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  files indexed
                </p>
              </div>
              {stats.languageBreakdown.length > 0 && (
                <div className="space-y-2">
                  {stats.languageBreakdown.slice(0, 6).map((lang) => (
                    <div key={lang.language} className="space-y-0.5">
                      <div className="flex items-center justify-between text-xs">
                        <span className="font-medium">
                          {lang.language || "Unknown"}
                        </span>
                        <span className="text-muted-foreground">
                          {lang.fileCount} ({Math.round(lang.percentage)}%)
                        </span>
                      </div>
                      <div className="h-1.5 w-full rounded-full bg-secondary">
                        <div
                          className="h-1.5 rounded-full bg-primary"
                          style={{ width: `${lang.percentage}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Last 5 sent newsletters */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Recent Newsletters
          </CardTitle>
          <Icon type="history" className="!text-base text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {newslettersResponse.isFetching ? (
            <div className="flex items-center gap-2 py-2">
              <Spinner className="h-4 w-4" />
              <span className="text-sm text-muted-foreground">
                Loading newsletters…
              </span>
            </div>
          ) : sentNewsletters.length === 0 ? (
            <EmptyState
              icon="mail"
              title="Nothing sent yet"
              description="Send your first newsletter to see it here."
              cta={
                <Link
                  href={`/repositories/${repoId}/newsletters`}
                  className="text-sm text-primary hover:underline"
                >
                  Go to Newsletters
                </Link>
              }
            />
          ) : (
            <ul className="divide-y divide-border">
              {sentNewsletters.map((newsletter) => (
                <li key={newsletter.id}>
                  <button
                    type="button"
                    className="flex w-full items-center justify-between gap-4 py-2.5 text-left hover:text-primary"
                    onClick={() =>
                      router.push(
                        `/repositories/${repoId}/newsletters/${newsletter.id}`,
                      )
                    }
                  >
                    <span className="min-w-0 flex-1 truncate text-sm font-medium">
                      {newsletter.subject || "Newsletter"}
                    </span>
                    <span className="flex shrink-0 items-center gap-3 text-xs text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Icon type="visibility" className="!text-sm" />
                        {Math.round(newsletter.openRate * 100)}%
                      </span>
                      <span className="flex items-center gap-1">
                        <Icon type="ads_click" className="!text-sm" />
                        {Math.round(newsletter.clickRate * 100)}%
                      </span>
                      <span>
                        {formatDaysAgo(newsletter.sentAt ?? "") ?? ""}
                      </span>
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
