// RichTextEditor uses Tiptap and several API clients. Mock heavy dependencies
// so the smoke test focuses only on render behaviour of our wrapper.
jest.mock("@/api/files", () => ({ default: { upload: jest.fn() } }));
jest.mock("@/lib/emoji-extension", () => ({ EmojiExtension: {} }));
jest.mock("@/lib/image-extension", () => ({ CustomImage: {} }));

import { render } from "@testing-library/react";
import { RichTextEditor } from "../rich-text-editor";

describe("RichTextEditor", () => {
  it("renders without crashing", () => {
    render(<RichTextEditor />);
  });

  it("renders in read-only mode without crashing", () => {
    render(<RichTextEditor editable={false} value="<p>Hello</p>" />);
  });
});
