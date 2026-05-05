/**
 * Returns the uppercased initials for a display name (max 2 characters).
 * "Jane Doe" → "JD", "Alice" → "A", null/undefined → fallback
 */
export function getInitials(name?: string | null, fallback = ""): string {
  if (!name) return fallback;
  return name
    .split(" ")
    .map((n) => n.charAt(0))
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

const AVATAR_COLORS: { backgroundColor: string; color: string }[] = [
  { backgroundColor: "#FDE8E8", color: "#9B2335" }, // soft rose / deep rose
  { backgroundColor: "#FEF3C7", color: "#92400E" }, // pale amber / dark amber
  { backgroundColor: "#D1FAE5", color: "#065F46" }, // mint green / dark green
  { backgroundColor: "#DBEAFE", color: "#1E40AF" }, // sky blue / dark blue
  { backgroundColor: "#EDE9FE", color: "#5B21B6" }, // lavender / deep purple
  { backgroundColor: "#FCE7F3", color: "#9D174D" }, // blush pink / dark pink
  { backgroundColor: "#CCFBF1", color: "#115E59" }, // aqua / dark teal
  { backgroundColor: "#FEE2E2", color: "#991B1B" }, // peach / dark red
  { backgroundColor: "#E0F2FE", color: "#0C4A6E" }, // light cyan / navy
  { backgroundColor: "#F3F4F6", color: "#374151" }, // cool grey / charcoal
];

/**
 * Returns a deterministic background + text color style for an avatar fallback.
 * The same name always maps to the same color by summing char codes mod 10.
 */
export function getAvatarColor(name?: string | null): {
  backgroundColor: string;
  color: string;
} {
  if (!name) return AVATAR_COLORS[0];
  const index =
    name.split("").reduce((sum, char) => sum + char.charCodeAt(0), 0) %
    AVATAR_COLORS.length;
  return AVATAR_COLORS[index];
}
