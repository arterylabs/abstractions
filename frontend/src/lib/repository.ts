import { isoWeekdayToName, utcDayHourToLocal } from "@/lib/timezone";
import { Repository } from "@/types";

/** "9:00 AM" for an hour of the day expressed as 0–23. */
export function formatHour(hour: number): string {
  return new Date(0, 0, 1, hour).toLocaleTimeString("en-US", {
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  });
}

/** "Monday at 9:00 AM UTC" — the send window as it is stored. */
export function formatSendWindow(sendDay: number, sendHour: number): string {
  return `${isoWeekdayToName(sendDay)} at ${formatHour(sendHour)} UTC`;
}

/** The same window translated into the reader's own timezone. */
export function formatLocalSendWindow(
  sendDay: number,
  sendHour: number,
  timezone: string,
): string {
  const { localDay, localHour } = utcDayHourToLocal(
    sendDay,
    sendHour,
    timezone,
  );
  return `${isoWeekdayToName(localDay)} at ${formatHour(localHour)} (${timezone})`;
}

/**
 * One line describing who an issue goes out to. Workspace members and custom
 * addresses are two independent switches, so all four combinations are read.
 */
export function describeRecipients(
  repository: Pick<Repository, "emails" | "sendToWorkspaceMembers">,
): string {
  const emailCount = repository.emails?.length ?? 0;
  const toMembers = repository.sendToWorkspaceMembers ?? true;

  if (toMembers && emailCount > 0) {
    return `Workspace + ${emailCount} email${emailCount === 1 ? "" : "s"}`;
  }
  if (toMembers) {
    return "All workspace members";
  }
  if (emailCount > 0) {
    return `${emailCount} email${emailCount === 1 ? "" : "s"}`;
  }
  return "No one";
}

/** The web URL for the connected repository, when we can derive one. */
export function getRepositoryURL(repository: Repository): string | null {
  if (repository.serviceType !== "github" || !repository.fullName) {
    return null;
  }
  return `https://github.com/${repository.fullName}`;
}
