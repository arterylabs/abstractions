import { format, formatDistanceToNow, isValid, parseISO } from "date-fns";

/**
 * Formats an ISO date string as a short date.
 * @param value - ISO date string or nullish
 * @param fmt   - date-fns format string (default: "d MMM yyyy" → "5 Apr 2025")
 */
export function formatDate(
  value?: string | null,
  fmt = "d MMM yyyy",
): string | null {
  if (!value) return null;
  const date = parseISO(value);
  if (!isValid(date)) return value;
  try {
    return format(date, fmt);
  } catch {
    return value;
  }
}

/**
 * Formats an ISO date string as a time string.
 * @param value - ISO date string or nullish
 * @param fmt   - date-fns format string (default: "h:mm a" → "3:45 PM")
 */
export function formatTime(
  value?: string | null,
  fmt = "h:mm a",
): string | null {
  if (!value) return null;
  const date = parseISO(value);
  if (!isValid(date)) return null;
  try {
    return format(date, fmt);
  } catch {
    return null;
  }
}

/**
 * Formats an ISO date string as a date + time string.
 * @param value - ISO date string or nullish
 * @param fmt   - date-fns format string (default: "d MMM yyyy, h:mm a" → "5 Apr 2025, 3:45 PM")
 */
export function formatDateTime(
  value?: string | null,
  fmt = "d MMM yyyy, h:mm a",
): string | null {
  if (!value) return null;
  const date = parseISO(value);
  if (!isValid(date)) return value;
  try {
    return format(date, fmt);
  } catch {
    return value;
  }
}

export function formatDaysAgo(value?: string | null): string | null {
  if (!value) return null;
  const date = parseISO(value);
  if (!isValid(date)) return value;
  try {
    return formatDistanceToNow(date, { addSuffix: true });
  } catch {
    return value;
  }
}

export function formatRate(rate: number): string {
  return `${(rate * 100).toFixed(1)}%`;
}
