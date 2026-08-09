"use client";

import { useEffect } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { formatDateTime } from "@/lib/format";
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
        <p className="press-label text-muted-foreground">Pulling the issue…</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        eyebrow="Back issue"
        title={
          formatDateTime(newsletter.sentAt ?? newsletter.createdAt) ?? "Unknown"
        }
        showBackButton
        handleBack={() => router.push(`/repositories/${repoId}/newsletters`)}
      />

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-[220px_1fr]">
        {/* Colophon rail — this issue's circulation figures */}
        <dl className="flex flex-row gap-8 lg:flex-col">
          <div className="flex flex-col gap-1.5">
            <dt className="press-label text-muted-foreground">Recipients</dt>
            <dd className="press-display text-4xl tabular-nums text-primary">
              {newsletter.recipientCount}
            </dd>
          </div>
        </dl>

        <NewsletterPreviewFrame htmlContent={newsletter.htmlContent} />
      </div>
    </div>
  );
}
