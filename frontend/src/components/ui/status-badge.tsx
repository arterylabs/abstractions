import { Badge } from "@/components/ui/badge";

type StatusVariant =
  "success" | "destructive" | "secondary" | "outline" | "warning" | "default";

const STATUS_CONFIG: Record<string, { label: string; variant: StatusVariant }> =
  {
    pending: { label: "Pending", variant: "warning" },
    accepted: { label: "Accepted", variant: "success" },
    declined: { label: "Declined", variant: "destructive" },
    maybe: { label: "Maybe", variant: "secondary" },
    expired: { label: "Expired", variant: "destructive" },
    ignored: { label: "Ignored", variant: "secondary" },
    active: { label: "Active", variant: "success" },
    inactive: { label: "Inactive", variant: "secondary" },
    queued: { label: "Queued", variant: "secondary" },
    sending: { label: "Sending", variant: "warning" },
    sent: { label: "Sent", variant: "warning" },
    delivered: { label: "Delivered", variant: "success" },
    scheduled: { label: "Scheduled", variant: "secondary" },
    undelivered: { label: "Undelivered", variant: "destructive" },
    canceled: { label: "Canceled", variant: "destructive" },
    success: { label: "Success", variant: "success" },
    failed: { label: "Failed", variant: "destructive" },
    partial_success: { label: "Partial Success", variant: "warning" },
  };

interface StatusBadgeProps {
  status: string;
  /** Override the display label. Defaults to the capitalized status key from config. */
  label?: string;
  className?: string;
}

export function StatusBadge({ status, label, className }: StatusBadgeProps) {
  const config = STATUS_CONFIG[status.toLowerCase()];
  return (
    <Badge variant={config?.variant ?? "secondary"} className={className}>
      {label ?? config?.label ?? status}
    </Badge>
  );
}
