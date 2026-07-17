"use client";

import React, {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useState,
} from "react";
import { cn } from "@/lib/utils";
import type { EmojiResult } from "@/lib/emoji-data";

export interface EmojiSuggestionListRef {
  onKeyDown: (event: KeyboardEvent) => boolean;
}

interface EmojiSuggestionListProps {
  items: EmojiResult[];
  /** Bounding rect of the cursor — used to position the popover. */
  rect: (() => DOMRect | null) | null;
  onSelect: (item: EmojiResult) => void;
}

export const EmojiSuggestionList = forwardRef<
  EmojiSuggestionListRef,
  EmojiSuggestionListProps
>(({ items, rect, onSelect }, ref) => {
  const [selectedIndex, setSelectedIndex] = useState(0);

  // Reset selection whenever the item list changes
  useEffect(() => setSelectedIndex(0), [items]);

  useImperativeHandle(ref, () => ({
    onKeyDown(event: KeyboardEvent) {
      if (event.key === "ArrowDown") {
        setSelectedIndex((i) => (i + 1) % Math.max(items.length, 1));
        return true;
      }
      if (event.key === "ArrowUp") {
        setSelectedIndex((i) =>
          i === 0 ? Math.max(items.length - 1, 0) : i - 1,
        );
        return true;
      }
      if (event.key === "Enter" || event.key === "Tab") {
        const item = items[selectedIndex];
        if (item) onSelect(item);
        return true;
      }
      return false;
    },
  }));

  if (!items.length) return null;

  const clientRect = rect?.();
  if (!clientRect) return null;

  const style: React.CSSProperties = {
    position: "fixed",
    top: clientRect.bottom + 6,
    left: clientRect.left,
    zIndex: 50,
  };

  return (
    <div
      style={style}
      className={cn(
        "min-w-[180px] max-h-[240px] overflow-y-auto",
        "rounded-md border border-border bg-popover shadow-md",
        "py-1 text-sm text-popover-foreground",
      )}
    >
      {items.map((item, index) => (
        <button
          key={item.id}
          type="button"
          className={cn(
            "flex w-full items-center gap-2 px-3 py-1.5 text-left",
            "hover:bg-accent hover:text-accent-foreground",
            index === selectedIndex && "bg-accent text-accent-foreground",
          )}
          onMouseEnter={() => setSelectedIndex(index)}
          onClick={() => onSelect(item)}
        >
          <span className="text-base leading-none">{item.emoji}</span>
          <span className="truncate text-xs text-muted-foreground">
            :{item.id}:
          </span>
        </button>
      ))}
    </div>
  );
});

EmojiSuggestionList.displayName = "EmojiSuggestionList";
