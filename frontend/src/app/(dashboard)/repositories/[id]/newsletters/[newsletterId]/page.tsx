"use client";

import { useEffect } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { formatDateTime } from "@/lib/format";
import { formatRate } from "@/lib/format";
import { useParams, useRouter } from "next/navigation";
import { useAppSelector } from "@/store/hooks";
import { NewsletterPreviewFrame } from "./NewsletterPreviewFrame";

export default function NewsletterDetailPage() {
  const router = useRouter();
  const params = useParams<{ id: string; newsletterId: string }>();
  const repoId = params.id;
  const newsletterId = params.newsletterId;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [newsletterResponse, loadNewsletter] = useAPI(api.newsletters.get);

  useEffect(() => {
    if (!repoId || !newsletterId || !currentWorkspace?.id) return;
    void loadNewsletter(currentWorkspace.id, repoId, newsletterId);
  }, [loadNewsletter, newsletterId, repoId, currentWorkspace?.id]);

  const newsletter = newsletterResponse.result;

  if (!newsletter || newsletterResponse.isFetching) {
    return (
      <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
        <Spinner />
        <p className="text-sm text-muted-foreground">Loading newsletter...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={`Newsletter — ${formatDateTime(newsletter.sentAt ?? newsletter.createdAt) ?? "Unknown"}`}
        showBackButton
        handleBack={() => router.push(`/repositories/${repoId}/newsletters`)}
      />

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Recipients</CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-semibold">
            {newsletter.recipientCount}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Opens</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            <p className="text-2xl font-semibold">{newsletter.openCount}</p>
            <p className="text-sm text-muted-foreground">
              {formatRate(newsletter.openRate)} open rate
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Clicks</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            <p className="text-2xl font-semibold">{newsletter.clickCount}</p>
            <p className="text-sm text-muted-foreground">
              {formatRate(newsletter.clickRate)} click rate
            </p>
          </CardContent>
        </Card>
      </div>

      <NewsletterPreviewFrame htmlContent={newsletter.htmlContent} />
    </div>
  );
}
