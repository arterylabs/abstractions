import { act, render, screen } from "@testing-library/react";
import { createRef } from "react";
import userEvent from "@testing-library/user-event";
import { EmojiSuggestionList } from "../emoji-suggestion-list";
import type { EmojiSuggestionListRef } from "../emoji-suggestion-list";
import type { EmojiResult } from "@/lib/emoji-data";

const mockRect = () =>
  ({
    top: 100,
    bottom: 120,
    left: 50,
    right: 200,
    width: 150,
    height: 20,
  }) as DOMRect;

const mockItems: EmojiResult[] = [
  { id: "smile", name: "Smile", emoji: "😊" },
  { id: "thumbsup", name: "Thumbs Up", emoji: "👍" },
];

describe("EmojiSuggestionList", () => {
  it("renders emoji items when rect is provided", () => {
    render(
      <EmojiSuggestionList
        items={mockItems}
        rect={mockRect}
        onSelect={jest.fn()}
      />,
    );
    expect(screen.getByText(":smile:")).toBeInTheDocument();
    expect(screen.getByText(":thumbsup:")).toBeInTheDocument();
  });

  it("renders nothing when items list is empty", () => {
    const { container } = render(
      <EmojiSuggestionList items={[]} rect={mockRect} onSelect={jest.fn()} />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it("renders nothing when rect is null", () => {
    const { container } = render(
      <EmojiSuggestionList
        items={mockItems}
        rect={null}
        onSelect={jest.fn()}
      />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it("onKeyDown ArrowDown advances the selected index", () => {
    const ref = createRef<EmojiSuggestionListRef>();
    render(
      <EmojiSuggestionList
        ref={ref}
        items={mockItems}
        rect={mockRect}
        onSelect={jest.fn()}
      />,
    );

    const buttons = screen.getAllByRole("button");
    // Initially the first item is highlighted
    expect(buttons[0]).toHaveClass("bg-accent");

    act(() => {
      ref.current?.onKeyDown(
        new KeyboardEvent("keydown", { key: "ArrowDown" }),
      );
    });

    // After ArrowDown the second item is highlighted
    expect(buttons[1]).toHaveClass("bg-accent");
  });

  it("onKeyDown ArrowUp wraps from first item to last", () => {
    const ref = createRef<EmojiSuggestionListRef>();
    render(
      <EmojiSuggestionList
        ref={ref}
        items={mockItems}
        rect={mockRect}
        onSelect={jest.fn()}
      />,
    );

    act(() => {
      ref.current?.onKeyDown(new KeyboardEvent("keydown", { key: "ArrowUp" }));
    });

    const buttons = screen.getAllByRole("button");
    // Wraps from index 0 → last index (1)
    expect(buttons[1]).toHaveClass("bg-accent");
  });

  it("onKeyDown Enter calls onSelect with the currently highlighted item", () => {
    const ref = createRef<EmojiSuggestionListRef>();
    const onSelect = jest.fn();
    render(
      <EmojiSuggestionList
        ref={ref}
        items={mockItems}
        rect={mockRect}
        onSelect={onSelect}
      />,
    );

    act(() => {
      ref.current?.onKeyDown(new KeyboardEvent("keydown", { key: "Enter" }));
    });

    // Default selectedIndex is 0 → first item
    expect(onSelect).toHaveBeenCalledWith(mockItems[0]);
  });

  it("onKeyDown Tab calls onSelect with the currently highlighted item", () => {
    const ref = createRef<EmojiSuggestionListRef>();
    const onSelect = jest.fn();
    render(
      <EmojiSuggestionList
        ref={ref}
        items={mockItems}
        rect={mockRect}
        onSelect={onSelect}
      />,
    );

    act(() => {
      ref.current?.onKeyDown(new KeyboardEvent("keydown", { key: "Tab" }));
    });

    expect(onSelect).toHaveBeenCalledWith(mockItems[0]);
  });

  it("onKeyDown returns false for unhandled keys", () => {
    const ref = createRef<EmojiSuggestionListRef>();
    render(
      <EmojiSuggestionList
        ref={ref}
        items={mockItems}
        rect={mockRect}
        onSelect={jest.fn()}
      />,
    );

    let result: boolean | undefined;
    act(() => {
      result = ref.current?.onKeyDown(
        new KeyboardEvent("keydown", { key: "Escape" }),
      );
    });

    expect(result).toBe(false);
  });

  it("clicking an emoji button calls onSelect with the corresponding item", async () => {
    const user = userEvent.setup();
    const onSelect = jest.fn();
    render(
      <EmojiSuggestionList
        items={mockItems}
        rect={mockRect}
        onSelect={onSelect}
      />,
    );

    // Click the second item's button
    const thumbsUpButton = screen
      .getByText(":thumbsup:")
      .closest("button") as HTMLButtonElement;
    await user.click(thumbsUpButton);

    expect(onSelect).toHaveBeenCalledWith(mockItems[1]);
  });
});
