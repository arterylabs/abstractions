"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { formatDateTime } from "@/lib/format";
import { getErrorMessage } from "@/lib/errors";
import { formatRate } from "@/lib/format";
import { useParams, useRouter } from "next/navigation";

export function NewsletterPreviewFrame({
  htmlContent,
}: {
  htmlContent: string;
}) {
  return (
    <iframe
      srcDoc={htmlContent}
      sandbox="allow-same-origin"
      title="Newsletter Preview"
      // colorScheme forces light mode inside the iframe regardless of OS preference
      style={{ colorScheme: "light" }}
      className="w-full min-h-[600px] border rounded-lg bg-white"
    />
  );
}

export default function NewsletterDetailPage() {
  const router = useRouter();
  const params = useParams<{ id: string; newsletterId: string }>();
  const repoId = params.id;
  const newsletterId = params.newsletterId;
  const [newsletterResponse, loadNewsletter] = useAPI(api.newsletters.get);
  const [emailContentResponse, loadEmailContent] = useAPI(
    api.newsletters.getEmailContent,
  );

  const [isEditMode, setIsEditMode] = useState(false);
  const [editedHtml, setEditedHtml] = useState("");
  const [isSavingHtml, setIsSavingHtml] = useState(false);
  const [isSending, setIsSending] = useState(false);

  // Debounced preview: update the iframe preview 300ms after user stops typing
  const [previewHtml, setPreviewHtml] = useState("");
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const baseHtml =
    emailContentResponse.result ?? newsletterResponse.result?.htmlContent ?? "";

  useEffect(() => {
    if (!repoId || !newsletterId) return;
    void loadNewsletter(repoId, newsletterId);
    void loadEmailContent(repoId, newsletterId);
  }, [loadEmailContent, loadNewsletter, newsletterId, repoId]);

  // Initialise editable HTML when content loads
  useEffect(() => {
    if (baseHtml && !editedHtml) {
      setEditedHtml(baseHtml);
      setPreviewHtml(baseHtml);
    }
  }, [baseHtml, editedHtml]);

  const handleHtmlChange = useCallback((value: string) => {
    setEditedHtml(value);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => setPreviewHtml(value), 300);
  }, []);

  const handleSaveHtml = async () => {
    setIsSavingHtml(true);
    try {
      await api.newsletters.updateContent(repoId, newsletterId, editedHtml);
      void loadEmailContent(repoId, newsletterId);
      Notifier.success({
        title: "Saved",
        message: "Newsletter content updated.",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSavingHtml(false);
    }
  };

  const handleSend = async () => {
    setIsSending(true);
    try {
      await api.newsletters.send(repoId, newsletterId);
      void loadNewsletter(repoId, newsletterId);
      Notifier.success({
        title: "Newsletter sent",
        message: "Your newsletter has been delivered to subscribers.",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to send",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSending(false);
    }
  };

  const newsletter = newsletterResponse.result;
  const isDraft =
    newsletter?.sentAt === null || newsletter?.sentAt === undefined;

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
        title={
          <span className="flex items-center gap-2">
            {`Newsletter — ${formatDateTime(newsletter.sentAt ?? newsletter.createdAt) ?? "Draft"}`}
            {isDraft && <Badge variant="secondary">Draft</Badge>}
          </span>
        }
        showBackButton
        handleBack={() => router.push(`/repositories/${repoId}/newsletters`)}
      >
        {isDraft && (
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => setIsEditMode((v) => !v)}
            >
              {isEditMode ? "Preview" : "Edit HTML"}
            </Button>
            {isEditMode && (
              <Button
                size="sm"
                variant="outline"
                onClick={() => void handleSaveHtml()}
                loading={isSavingHtml}
              >
                Save
              </Button>
            )}
            <Button
              size="sm"
              onClick={() => void handleSend()}
              loading={isSending}
            >
              Send
            </Button>
          </div>
        )}
      </PageHeader>

      {!isDraft && (
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
      )}

      {isEditMode ? (
        <textarea
          className="w-full min-h-[600px] rounded-lg border bg-background p-4 font-mono text-xs focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-y"
          value={editedHtml}
          onChange={(e) => handleHtmlChange(e.target.value)}
          spellCheck={false}
        />
      ) : (
        <NewsletterPreviewFrame htmlContent={previewHtml || baseHtml} />
      )}
    </div>
  );
}
