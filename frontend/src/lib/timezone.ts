export function getBrowserTimezone(): string {
  return Intl.DateTimeFormat().resolvedOptions().timeZone;
}

// Returns the UTC offset for a timezone in fractional hours (e.g. 5.5 for IST).
// Uses the current moment so DST is accounted for.
function getTimezoneOffsetHours(timezone: string): number {
  const now = new Date();
  const formatter = new Intl.DateTimeFormat("en-US", {
    timeZone: timezone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
  const parts = formatter.formatToParts(now);
  const get = (type: string) =>
    parseInt(parts.find((p) => p.type === type)?.value ?? "0");

  const tzHour = get("hour") % 24;
  const tzMinute = get("minute");
  const tzYear = get("year");
  const tzMonth = get("month") - 1;
  const tzDay = get("day");

  const tzMs = Date.UTC(tzYear, tzMonth, tzDay, tzHour, tzMinute);
  const utcMs = Date.UTC(
    now.getUTCFullYear(),
    now.getUTCMonth(),
    now.getUTCDate(),
    now.getUTCHours(),
    now.getUTCMinutes(),
  );

  let offsetMinutes = (tzMs - utcMs) / (1000 * 60);
  if (offsetMinutes > 720) offsetMinutes -= 1440;
  if (offsetMinutes < -720) offsetMinutes += 1440;

  return offsetMinutes / 60;
}

// Converts stored UTC day+hour back to local values for display in the given timezone.
// Note: sub-hour offsets (e.g. IST UTC+5:30) are rounded to the nearest whole hour.
export function utcDayHourToLocal(
  utcDay: number,
  utcHour: number,
  timezone: string,
): { localDay: number; localHour: number } {
  const offsetHours = getTimezoneOffsetHours(timezone);
  const localHourRaw = utcHour + offsetHours;
  let localHour = Math.round(localHourRaw);
  let dayDelta = 0;

  if (localHour < 0) {
    localHour += 24;
    dayDelta = -1;
  } else if (localHour >= 24) {
    localHour -= 24;
    dayDelta = 1;
  }

  const localDay = ((utcDay - 1 + dayDelta + 7) % 7) + 1;
  return { localDay, localHour };
}

export function isoWeekdayToName(day: number): string {
  const dayNames = [
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
    "Sunday",
  ];

  return dayNames[day - 1] ?? "Unknown";
}

export function getDayOptions(): Array<{ value: number; label: string }> {
  return Array.from({ length: 7 }, (_, index) => ({
    value: index + 1,
    label: isoWeekdayToName(index + 1),
  }));
}

export function daysUntilNextSend(
  sendDay: number,
  sendHour: number,
  now = new Date(),
): number {
  const currentUtcDay = ((now.getUTCDay() + 6) % 7) + 1;
  const dayOffset = (sendDay - currentUtcDay + 7) % 7;
  const nextSendAt = new Date(
    Date.UTC(
      now.getUTCFullYear(),
      now.getUTCMonth(),
      now.getUTCDate(),
      sendHour,
    ),
  );

  nextSendAt.setUTCDate(nextSendAt.getUTCDate() + dayOffset);

  if (dayOffset === 0 && nextSendAt <= now) {
    nextSendAt.setUTCDate(nextSendAt.getUTCDate() + 7);
  }

  const startOfToday = Date.UTC(
    now.getUTCFullYear(),
    now.getUTCMonth(),
    now.getUTCDate(),
  );
  const startOfSendDay = Date.UTC(
    nextSendAt.getUTCFullYear(),
    nextSendAt.getUTCMonth(),
    nextSendAt.getUTCDate(),
  );

  return Math.round((startOfSendDay - startOfToday) / (1000 * 60 * 60 * 24));
}
