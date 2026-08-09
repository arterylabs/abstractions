import * as React from "react";
import { cn } from "@/lib/utils";

interface SectionMarkProps extends React.ComponentProps<"div"> {
  /** Running head, e.g. "Delivery schedule". */
  label: React.ReactNode;
  /** Optional folio number set before the label, e.g. "01". */
  index?: string;
  /** Actions pinned to the far end of the rule. */
  action?: React.ReactNode;
}

/**
 * Opens a section of a page the way a folio opens a section of the paper: an
 * uppercase mono running head, optionally numbered, with a hairline rule
 * carrying the eye across to any actions on the right.
 *
 * This is the in-app counterpart of `SectionMark` on the landing page.
 */
function SectionMark({
  label,
  index,
  action,
  className,
  ...props
}: SectionMarkProps) {
  return (
    <div className={cn("flex items-center gap-4", className)} {...props}>
      {index && <span className="press-label text-primary">§ {index}</span>}
      <span className="press-label shrink-0 text-muted-foreground">
        {label}
      </span>
      <span aria-hidden className="h-px flex-1 bg-border" />
      {action && <div className="shrink-0">{action}</div>}
    </div>
  );
}

export { SectionMark };
