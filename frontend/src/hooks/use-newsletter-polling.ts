import { useEffect } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";

// Statuses that indicate a newsletter has finished processing.
const TERMINAL_STATUSES = new Set(["sent", "error"]);

// How often to re-fetch the list while at least one newsletter is still in
// progress. 5 s balances UI freshness against server load.
const POLL_INTERVAL_MS = 5000;

/**
 * Fetches the newsletter list for a repository and automatically polls every
 * 5 seconds while any newsletter is in a non-terminal status (i.e. "queued"
 * or "processing"). Polling stops once all newsletters are terminal or the
 * component using this hook unmounts.
 *
 * Returns the same [response, refresh] tuple as useAPI so callers can still
 * trigger manual refreshes (e.g. after a resend action).
 */
export function useNewsletterPolling(
  workspaceId: string | undefined,
  repoId: string,
) {
  const [response, load] = useAPI(api.newsletters.list);

  // Initial fetch
  useEffect(() => {
    if (!workspaceId || !repoId) return;
    void load(workspaceId, repoId);
  }, [load, workspaceId, repoId]);

  // Polling — only active when at least one newsletter is non-terminal
  useEffect(() => {
    if (!workspaceId || !repoId) return;
    const hasNonTerminal = response.result?.some(
      (n) => !TERMINAL_STATUSES.has(n.status),
    );
    if (!hasNonTerminal) return;

    const id = window.setInterval(
      () => void load(workspaceId, repoId),
      POLL_INTERVAL_MS,
    );
    return () => window.clearInterval(id);
  }, [response.result, load, workspaceId, repoId]);

  return [response, load] as const;
}
