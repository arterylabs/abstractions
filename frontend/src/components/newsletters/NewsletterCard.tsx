"use client";

import type { NewsletterSummary } from "@/api/newsletters";
import {
  Card,
  CardHeader,
  CardDescription,
  CardTitle,
} from "@/components/ui/card";
import Tooltip from "@/components/ui/tooltip";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { SimpleDropdownMenu } from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";
import { formatDateTime, formatDaysAgo } from "@/lib/format";

interface NewsletterCardProps {
  newsletter: NewsletterSummary;
  onClick: () => void;
  onResend?: () => void;
}

type BadgeVariant =
  | "info"
  | "warning"
  | "success"
  | "destructive"
  | "secondary"
  | "default"
  | "outline";

const STATUS_BADGE: Record<string, { label: string; variant: BadgeVariant }> = {
  queued: { label: "Queued", variant: "info" },
  processing: { label: "Processing", variant: "warning" },
  sent: { label: "Sent", variant: "success" },
  error: { label: "Error", variant: "destructive" },
};

export function NewsletterCard({
  newsletter,
  onClick,
  onResend,
}: NewsletterCardProps) {
  const sentAt = newsletter.sentAt ?? newsletter.createdAt;
  const badge = STATUS_BADGE[newsletter.status] ?? {
    label: newsletter.status,
    variant: "secondary" as BadgeVariant,
  };
  const isError = newsletter.status === "error";

  return (
    <div className="group relative">
      {/* Clickable card area */}
      <button type="button" className="w-full text-left" onClick={onClick}>
        <Card className="press-card overflow-hidden">
          <CardHeader className="gap-2.5 p-5">
            <CardTitle className="text-[1.6rem] transition-colors group-hover:text-primary">
              {newsletter.subject || "Newsletter"}
            </CardTitle>
            <CardDescription className="flex items-center gap-3 font-mono text-xs">
              <Tooltip value={formatDateTime(sentAt) ?? ""}>
                <span className="cursor-default">
                  {formatDaysAgo(sentAt) ?? "Draft"}
                </span>
              </Tooltip>
              {isError && newsletter.errorMessage ? (
                <Tooltip value={newsletter.errorMessage}>
                  <Badge
                    variant={badge.variant}
                    size="sm"
                    className="cursor-default"
                  >
                    {badge.label}
                  </Badge>
                </Tooltip>
              ) : (
                <Badge variant={badge.variant} size="sm">
                  {badge.label}
                </Badge>
              )}
            </CardDescription>
          </CardHeader>
        </Card>
      </button>

      {/* More actions — absolutely positioned so it sits outside the card button */}
      {isError && onResend && (
        <div className="absolute right-4 top-1/2 -translate-y-1/2">
          <SimpleDropdownMenu
            widthClassName="w-40"
            options={[
              {
                label: "Resend",
                icon: <Icon type="refresh" className="text-sm" />,
                onClick: onResend,
              },
            ]}
          >
            <Button variant="ghost" size="icon" aria-label="More actions">
              <Icon type="more_horiz" className="text-base" />
            </Button>
          </SimpleDropdownMenu>
        </div>
      )}
    </div>
  );
}
