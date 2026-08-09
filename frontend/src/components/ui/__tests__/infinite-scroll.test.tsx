import { fireEvent, render, screen } from "@testing-library/react";
import { InfiniteScroll } from "../infinite-scroll";

/** Put the feed near the bottom so handleScroll's threshold is met. */
function mockNearBottom(element: HTMLElement) {
  Object.defineProperty(element, "scrollTop", {
    configurable: true,
    value: 100,
  });
  Object.defineProperty(element, "scrollHeight", {
    configurable: true,
    value: 200,
  });
  Object.defineProperty(element, "clientHeight", {
    configurable: true,
    value: 100,
  });
}

describe("InfiniteScroll", () => {
  const defaultProps = {
    loading: false,
    hasMore: true,
    onLoadMore: jest.fn(),
  };

  beforeEach(() => {
    defaultProps.onLoadMore = jest.fn();
  });

  it("should render children", () => {
    render(
      <InfiniteScroll {...defaultProps}>
        <p>Item 1</p>
      </InfiniteScroll>,
    );
    expect(screen.getByText("Item 1")).toBeInTheDocument();
  });

  it("should show loading spinner when loading=true", () => {
    const { container } = render(
      <InfiniteScroll {...defaultProps} loading={true}>
        <p>Item 1</p>
      </InfiniteScroll>,
    );
    expect(container.querySelector(".animate-spin")).toBeInTheDocument();
  });

  it("should have aria-busy='true' when loading", () => {
    render(
      <InfiniteScroll {...defaultProps} loading={true}>
        <p>Item</p>
      </InfiniteScroll>,
    );
    expect(screen.getByRole("feed")).toHaveAttribute("aria-busy", "true");
  });

  it("should have role='feed'", () => {
    render(
      <InfiniteScroll {...defaultProps}>
        <p>Item</p>
      </InfiniteScroll>,
    );
    expect(screen.getByRole("feed")).toBeInTheDocument();
  });

  it("should have aria-label='Infinite scroll content'", () => {
    render(
      <InfiniteScroll {...defaultProps}>
        <p>Item</p>
      </InfiniteScroll>,
    );
    expect(
      screen.getByRole("feed", { name: "Infinite scroll content" }),
    ).toBeInTheDocument();
  });

  it("should call onLoadMore when scrolled near the bottom", () => {
    render(
      <InfiniteScroll {...defaultProps}>
        <p>Item 1</p>
      </InfiniteScroll>,
    );

    const feed = screen.getByRole("feed");
    mockNearBottom(feed);
    fireEvent.scroll(feed);

    expect(defaultProps.onLoadMore).toHaveBeenCalledTimes(1);
  });

  it("should not call onLoadMore when loading", () => {
    render(
      <InfiniteScroll {...defaultProps} loading={true}>
        <p>Item 1</p>
      </InfiniteScroll>,
    );

    const feed = screen.getByRole("feed");
    mockNearBottom(feed);
    fireEvent.scroll(feed);

    expect(defaultProps.onLoadMore).not.toHaveBeenCalled();
  });

  it("should not call onLoadMore when hasMore is false", () => {
    render(
      <InfiniteScroll {...defaultProps} hasMore={false}>
        <p>Item 1</p>
      </InfiniteScroll>,
    );

    const feed = screen.getByRole("feed");
    mockNearBottom(feed);
    fireEvent.scroll(feed);

    expect(defaultProps.onLoadMore).not.toHaveBeenCalled();
  });
});
