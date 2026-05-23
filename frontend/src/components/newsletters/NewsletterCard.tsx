"use client";

import type { NewsletterSummary } from "@/api/newsletters";
import {
  Card,
  CardHeader,
  CardDescription,
  CardTitle,
} from "@/components/ui/card";
import Tooltip from "@/components/ui/tooltip";
import { formatDateTime, formatDaysAgo } from "@/lib/format";

interface NewsletterCardProps {
  newsletter: NewsletterSummary;
  onClick: () => void;
}

export function NewsletterCard({ newsletter, onClick }: NewsletterCardProps) {
  const sentAt = newsletter.sentAt ?? newsletter.createdAt;

  return (
    <button type="button" className="group text-left" onClick={onClick}>
      <Card className="transition-colors hover:bg-accent">
        <CardHeader className="gap-1">
          <CardTitle className="text-base transition-colors group-hover:text-primary">
            {newsletter.subject || "Newsletter"}
          </CardTitle>
          <CardDescription>
            <Tooltip value={formatDateTime(sentAt) ?? ""}>
              <span className="cursor-default">
                {formatDaysAgo(sentAt) ?? "Draft"}
              </span>
            </Tooltip>
          </CardDescription>
        </CardHeader>
      </Card>
    </button>
  );
}
