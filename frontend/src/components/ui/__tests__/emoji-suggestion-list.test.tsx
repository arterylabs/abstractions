import { render, screen } from "@testing-library/react";
import { EmojiSuggestionList } from "../emoji-suggestion-list";
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
});
