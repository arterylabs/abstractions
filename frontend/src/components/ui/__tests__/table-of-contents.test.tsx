import { render, screen } from "@testing-library/react";
import TableOfContents, { type TocEntry } from "../table-of-contents";

// Mock IntersectionObserver before tests
const mockObserve = jest.fn();
const mockDisconnect = jest.fn();

class MockIntersectionObserver {
  observe = mockObserve;
  disconnect = mockDisconnect;
  unobserve = jest.fn();
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  constructor(_callback: IntersectionObserverCallback) {}
}

beforeAll(() => {
  Object.defineProperty(window, "IntersectionObserver", {
    writable: true,
    configurable: true,
    value: MockIntersectionObserver,
  });
});

beforeEach(() => {
  jest.clearAllMocks();
});

const toc: TocEntry[] = [
  { id: "intro", label: "Introduction", level: 2 },
  { id: "details", label: "Details", level: 2 },
  { id: "sub-details", label: "Sub Details", level: 3 },
];

describe("TableOfContents", () => {
  it("renders nothing when toc is empty", () => {
    const { container } = render(<TableOfContents toc={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("renders 'On this page' heading when toc has entries", () => {
    render(<TableOfContents toc={toc} />);
    expect(screen.getByText("On this page")).toBeInTheDocument();
  });

  it("renders all toc entry labels", () => {
    render(<TableOfContents toc={toc} />);
    expect(screen.getByText("Introduction")).toBeInTheDocument();
    expect(screen.getByText("Details")).toBeInTheDocument();
    expect(screen.getByText("Sub Details")).toBeInTheDocument();
  });

  it("renders links with correct href anchors", () => {
    render(<TableOfContents toc={toc} />);
    expect(screen.getByRole("link", { name: "Introduction" })).toHaveAttribute(
      "href",
      "#intro",
    );
    expect(screen.getByRole("link", { name: "Details" })).toHaveAttribute(
      "href",
      "#details",
    );
  });

  it("renders nav element", () => {
    render(<TableOfContents toc={toc} />);
    expect(screen.getByRole("navigation")).toBeInTheDocument();
  });
});
