import { Extension, Editor } from "@tiptap/core";
import Suggestion, {
  SuggestionProps,
  SuggestionKeyDownProps,
} from "@tiptap/suggestion";
import React, { createRef } from "react";
import { createRoot, Root } from "react-dom/client";
import { searchEmojis, EmojiResult } from "@/lib/emoji-data";
import {
  EmojiSuggestionList,
  EmojiSuggestionListRef,
} from "@/components/ui/emoji-suggestion-list";

/**
 * Tiptap extension that shows an emoji autocomplete popover when the user
 * types `:` followed by a search query. Selecting an emoji inserts its
 * Unicode character and removes the `:query` text.
 */
export const EmojiExtension = Extension.create({
  name: "emojiSuggestion",

  addProseMirrorPlugins() {
    return [
      Suggestion<EmojiResult>({
        editor: this.editor,
        char: ":",
        allowSpaces: false,

        // Fetch emoji matches for the current query
        items: async ({ query }: { query: string }) => {
          return searchEmojis(query);
        },

        // Insert the selected emoji Unicode character, replacing the :query range
        command: ({
          editor,
          range,
          props,
        }: {
          editor: Editor;
          range: { from: number; to: number };
          props: EmojiResult;
        }) => {
          editor
            .chain()
            .focus()
            .deleteRange(range)
            .insertContent(props.emoji)
            .run();
        },

        // Manage the React-rendered suggestion popover
        render: () => {
          let container: HTMLDivElement | null = null;
          let root: Root | null = null;
          const listRef = createRef<EmojiSuggestionListRef>();

          const renderComponent = (props: SuggestionProps<EmojiResult>) => {
            if (!root) return;
            root.render(
              React.createElement(EmojiSuggestionList, {
                ref: listRef,
                items: props.items,
                rect: props.clientRect ?? null,
                onSelect: (item: EmojiResult) => props.command(item),
              }),
            );
          };

          return {
            onStart(props: SuggestionProps<EmojiResult>) {
              container = document.createElement("div");
              document.body.appendChild(container);
              root = createRoot(container);
              renderComponent(props);
            },

            onUpdate(props: SuggestionProps<EmojiResult>) {
              renderComponent(props);
            },

            onKeyDown({ event }: SuggestionKeyDownProps) {
              // Intercept Escape to close the popover without inserting
              if (event.key === "Escape") {
                root?.unmount();
                container?.remove();
                root = null;
                container = null;
                return true;
              }
              return listRef.current?.onKeyDown(event) ?? false;
            },

            onExit() {
              root?.unmount();
              container?.remove();
              root = null;
              container = null;
            },
          };
        },
      }),
    ];
  },
});
