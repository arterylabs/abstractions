import Icon from "@/components/ui/icon";

/**
 * Persistent strip that signals the operator is in platform admin mode.
 * Uses the violet brand token so it stays distinct from product primary (cyan).
 */
export function AdminModeBanner() {
  return (
    <div
      className="flex shrink-0 items-center gap-2 border-b border-violet/25 bg-violet/10 px-4 py-2 text-sm text-violet-light md:px-8"
      role="status"
      aria-label="Admin mode active"
    >
      <Icon
        type="admin_panel_settings"
        className="text-base text-violet-light"
      />
      <span className="font-medium text-foreground">Admin mode</span>
      <span className="hidden text-muted-foreground sm:inline">
        — you can manage users, workspaces, and platform data
      </span>
    </div>
  );
}
