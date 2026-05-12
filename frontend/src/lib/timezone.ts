const getReferenceDateParts = () => {
  const now = new Date();
  return {
    year: now.getFullYear(),
    month: now.getMonth(),
    day: now.getDate(),
  };
};

export function utcHourToLocal(utcHour: number): string {
  const { year, month, day } = getReferenceDateParts();
  const date = new Date(Date.UTC(year, month, day, utcHour));

  return date.toLocaleTimeString([], {
    hour: "numeric",
    minute: "2-digit",
  });
}

export function localHourToUtc(localHour: number): number {
  const localOffsetFromUtc = -new Date().getTimezoneOffset() / 60;
  const utcHour = Math.round(localHour - localOffsetFromUtc);

  return ((utcHour % 24) + 24) % 24;
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
