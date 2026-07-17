import { Notifier } from "@/components/ui/sonner";

export async function copyToClipboard(text: string, type: string) {
  try {
    await navigator.clipboard.writeText(text);
    Notifier.success(`${type} copied to clipboard`);
  } catch (err) {
    console.error("Failed to copy: ", err);
    Notifier.error(`Failed to copy ${type}`);
  }
}
