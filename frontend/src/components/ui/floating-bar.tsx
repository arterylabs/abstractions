import { type ReactNode } from "react";
import { cn } from "@/lib/utils";

interface FloatingBarProps {
  /** Left-side content: icon + descriptive text */
  content: ReactNode;
  /** Primary action button (e.g. Save / Create) */
  saveButton: ReactNode;
  /** Secondary action button (e.g. Cancel / Back) */
  cancelButton: ReactNode;
  /** Controls visibility with a slide-up animation. Defaults to true. */
  visible?: boolean;
}

/**
 * A fixed bottom-center bar used for contextual actions such as saving unsaved
 * changes or submitting a new record. Slides up from the bottom when visible.
 */
export function FloatingBar({
  content,
  saveButton,
  cancelButton,
  visible = true,
}: FloatingBarProps) {
  return (
    <div
      className={cn(
        "fixed bottom-6 left-1/2 -translate-x-1/2 z-50",
        "flex min-w-[700px] items-center gap-3 rounded-sm border border-primary/25 bg-primary px-8 py-3 shadow-2xl",
        "transition-all duration-300 ease-out",
        visible
          ? "translate-y-0 opacity-100 pointer-events-auto"
          : "translate-y-4 opacity-0 pointer-events-none",
      )}
    >
      <div className="flex items-center gap-3 flex-1 min-w-0">{content}</div>
      <div className="flex items-center gap-2 ml-2 shrink-0">
        {cancelButton}
        {saveButton}
      </div>
    </div>
  );
}
