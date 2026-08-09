// RichTextEditor uses Tiptap and API clients. Mock heavy deps so tests focus
// on our wrapper's mode/toolbar behaviour.
jest.mock("@/api/files", () => ({ default: { uploadFile: jest.fn() } }));
jest.mock("@/lib/emoji-extension", () => ({ EmojiExtension: {} }));
jest.mock("@/lib/image-extension", () => ({ CustomImage: {} }));
jest.mock("@tiptap/starter-kit", () => ({
  __esModule: true,
  default: { configure: jest.fn(() => ({})) },
}));
jest.mock("@tiptap/extension-placeholder", () => ({
  __esModule: true,
  default: { configure: jest.fn(() => ({})) },
}));
jest.mock("@tiptap/pm/state", () => ({
  TextSelection: { atStart: jest.fn(() => ({})) },
}));

const chainMock = {
  focus: jest.fn().mockReturnThis(),
  toggleBold: jest.fn().mockReturnThis(),
  toggleItalic: jest.fn().mockReturnThis(),
  toggleHeading: jest.fn().mockReturnThis(),
  toggleBulletList: jest.fn().mockReturnThis(),
  toggleOrderedList: jest.fn().mockReturnThis(),
  toggleBlockquote: jest.fn().mockReturnThis(),
  setImage: jest.fn().mockReturnThis(),
  insertContent: jest.fn().mockReturnThis(),
  run: jest.fn(),
};

function createMockEditor() {
  return {
    isActive: jest.fn(() => false),
    isEditable: true,
    getHTML: jest.fn(() => "<p></p>"),
    chain: jest.fn(() => chainMock),
    commands: {
      setContent: jest.fn(),
      insertContent: jest.fn(),
    },
    setEditable: jest.fn(),
    state: {
      doc: { descendants: jest.fn() },
      tr: {
        setSelection: jest.fn().mockReturnThis(),
        setNodeMarkup: jest.fn().mockReturnThis(),
        delete: jest.fn().mockReturnThis(),
      },
    },
    view: {
      dispatch: jest.fn(),
      dom: { blur: jest.fn() },
    },
  };
}

let mockEditor: ReturnType<typeof createMockEditor> | null;

jest.mock("@tiptap/react", () => ({
  useEditor: jest.fn(() => mockEditor),
  EditorContent: () => <div data-testid="editor-content" />,
}));

import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RichTextEditor } from "../rich-text-editor";

describe("RichTextEditor", () => {
  beforeEach(() => {
    mockEditor = createMockEditor();
    jest.clearAllMocks();
  });

  it("renders without crashing", () => {
    render(<RichTextEditor />);
    expect(screen.getByTestId("editor-content")).toBeInTheDocument();
  });

  it("renders in read-only mode without crashing", () => {
    render(<RichTextEditor editable={false} value="<p>Hello</p>" />);
    expect(screen.getByTestId("editor-content")).toBeInTheDocument();
  });

  it("shows the formatting toolbar when editable", () => {
    render(<RichTextEditor editable={true} />);
    expect(screen.getByText("format_bold")).toBeInTheDocument();
    expect(screen.getByText("format_italic")).toBeInTheDocument();
    expect(screen.getByText("image")).toBeInTheDocument();
  });

  it("hides the toolbar when editable is false", () => {
    render(<RichTextEditor editable={false} value="<p>Hello</p>" />);
    expect(screen.queryByText("format_bold")).not.toBeInTheDocument();
    expect(screen.queryByText("image")).not.toBeInTheDocument();
  });

  it("hides the toolbar when formatting is disabled", () => {
    render(<RichTextEditor formattingDisabled />);
    expect(screen.queryByText("format_bold")).not.toBeInTheDocument();
    expect(screen.getByTestId("editor-content")).toBeInTheDocument();
  });

  it("returns null when the editor is not ready", () => {
    mockEditor = null;
    const { container } = render(<RichTextEditor />);
    expect(container).toBeEmptyDOMElement();
  });

  it("invokes editor commands when toolbar buttons are clicked", async () => {
    const user = userEvent.setup();
    render(<RichTextEditor />);

    await user.click(screen.getByText("format_bold").closest("button")!);
    expect(chainMock.toggleBold).toHaveBeenCalled();

    await user.click(screen.getByText("format_italic").closest("button")!);
    expect(chainMock.toggleItalic).toHaveBeenCalled();

    await user.click(screen.getByText("format_h2").closest("button")!);
    expect(chainMock.toggleHeading).toHaveBeenCalledWith({ level: 2 });

    await user.click(screen.getByText("format_h3").closest("button")!);
    expect(chainMock.toggleHeading).toHaveBeenCalledWith({ level: 3 });

    await user.click(
      screen.getByText("format_list_bulleted").closest("button")!,
    );
    expect(chainMock.toggleBulletList).toHaveBeenCalled();

    await user.click(
      screen.getByText("format_list_numbered").closest("button")!,
    );
    expect(chainMock.toggleOrderedList).toHaveBeenCalled();

    await user.click(screen.getByText("format_quote").closest("button")!);
    expect(chainMock.toggleBlockquote).toHaveBeenCalled();
  });
});
