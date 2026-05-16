"use client";

import type { NewsletterTopic } from "@/api/newsletters";
import type { Repository } from "@/types/repository";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Spinner from "@/components/ui/spinner";
import { daysUntilNextSend } from "@/lib/timezone";

function formatDaysUntil(days: number): string {
  if (days === 0) return "today";
  if (days === 1) return "tomorrow";
  return `in ${days} days`;
}

interface NextUpCardProps {
  nextTopic: NewsletterTopic | undefined;
  repository: Repository;
  isLoading: boolean;
  onSendNow?: () => void;
  isSending?: boolean;
  onPreview?: () => void;
  isGeneratingPreview?: boolean;
}

export function NextUpCard({
  nextTopic,
  repository,
  isLoading,
  onSendNow,
  isSending,
  onPreview,
  isGeneratingPreview,
}: NextUpCardProps) {
  const sendsIn = repository.newsletterEnabled
    ? formatDaysUntil(
        daysUntilNextSend(repository.sendDay, repository.sendHour),
      )
    : null;

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between gap-3">
        <CardTitle className="text-base">Next up</CardTitle>
        {nextTopic && (
          <div className="flex items-center gap-2">
            {onPreview && (
              <Button
                size="sm"
                variant="default"
                onClick={onPreview}
                loading={isGeneratingPreview}
              >
                Preview
              </Button>
            )}
            {onSendNow && (
              <Button
                size="sm"
                variant="outline"
                onClick={onSendNow}
                loading={isSending}
              >
                Send Now
              </Button>
            )}
          </div>
        )}
      </CardHeader>
      <CardContent className="space-y-1">
        {isLoading ? (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Spinner className="h-4 w-4" />
          </div>
        ) : nextTopic ? (
          <>
            <p className="text-sm font-medium">{nextTopic.title}</p>
            {nextTopic.description && (
              <p className="text-xs text-muted-foreground">
                {nextTopic.description}
              </p>
            )}
            {sendsIn ? (
              <p className="text-xs text-muted-foreground">
                Scheduled to send {sendsIn}
              </p>
            ) : (
              <p className="text-xs text-muted-foreground italic">
                Newsletters are disabled — no email will be sent
              </p>
            )}
          </>
        ) : (
          <p className="text-sm text-muted-foreground italic">
            No topics queued yet
          </p>
        )}
      </CardContent>
    </Card>
  );
}
