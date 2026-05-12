"use client";

import type { NewsletterSummary } from "@/api/newsletters";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Icon from "@/components/ui/icon";
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
        <CardContent className="flex items-center gap-4 text-xs text-muted-foreground">
          <span className="flex items-center gap-1">
            <Icon type="visibility" className="text-sm" />
            {newsletter.openCount} opens
          </span>
          <span className="flex items-center gap-1">
            <Icon type="ads_click" className="text-sm" />
            {newsletter.clickCount} clicks
          </span>
        </CardContent>
      </Card>
    </button>
  );
}
