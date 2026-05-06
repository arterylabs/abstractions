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
