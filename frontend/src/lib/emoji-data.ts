import type { EmojiMartData } from "@emoji-mart/data";

export interface EmojiResult {
  id: string;
  name: string;
  emoji: string;
}

let cachedEmojis: EmojiResult[] | null = null;

async function getAllEmojis(): Promise<EmojiResult[]> {
  if (cachedEmojis) return cachedEmojis;

  // Dynamically import to keep it out of the initial bundle
  const data = (await import("@emoji-mart/data")).default as EmojiMartData;

  cachedEmojis = Object.values(data.emojis)
    .filter((e) => e.skins?.[0]?.native)
    .map((e) => ({
      id: e.id,
      name: e.name,
      emoji: e.skins[0].native,
    }));

  return cachedEmojis;
}

/**
 * Search emojis by short code or name.
 * Returns up to 10 results. If query is empty, returns a popular starter set.
 */
export async function searchEmojis(query: string): Promise<EmojiResult[]> {
  const all = await getAllEmojis();

  if (!query) {
    // Default set of popular emojis shown before the user types anything
    const popularIds = [
      "grinning",
      "smile",
      "heart_eyes",
      "joy",
      "fire",
      "thumbsup",
      "+1",
      "clap",
      "tada",
      "rocket",
    ];
    const popular = popularIds
      .map((id) => all.find((e) => e.id === id))
      .filter((e): e is EmojiResult => !!e);
    return popular.slice(0, 10);
  }

  const q = query.toLowerCase();
  return all
    .filter((e) => e.id.includes(q) || e.name.toLowerCase().includes(q))
    .slice(0, 10);
}
