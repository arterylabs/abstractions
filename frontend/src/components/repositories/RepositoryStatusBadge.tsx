import { Badge } from "@/components/ui/badge";
import Icon from "@/components/ui/icon";
import Tooltip from "@/components/ui/tooltip";
import { formatDateTime } from "@/lib/format";
import { Repository } from "@/types";

const STATUS_CONFIG: Record<
  Repository["status"],
  {
    label: string;
    variant: "secondary" | "warning" | "success" | "destructive";
    className?: string;
  }
> = {
  pending: { label: "Pending", variant: "secondary" },
  indexing: {
    label: "Indexing…",
    variant: "warning",
    className: "animate-pulse",
  },
  indexed: { label: "Indexed", variant: "success" },
  error: { label: "Error", variant: "destructive" },
};

export function RepositoryStatusBadge({
  status,
  indexedAt,
}: {
  status: Repository["status"];
  indexedAt?: string | null;
}) {
  const config = STATUS_CONFIG[status];
  const badge = (
    <Badge variant={config.variant} className={config.className}>
      <Icon type="database" className="text-[10px]" />
      {config.label}
    </Badge>
  );
  if (status === "indexed" && indexedAt) {
    return (
      <Tooltip value={`Last indexed ${formatDateTime(indexedAt)}`}>
        {badge}
      </Tooltip>
    );
  }
  return badge;
}
