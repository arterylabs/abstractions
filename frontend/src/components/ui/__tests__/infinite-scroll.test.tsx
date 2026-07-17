import { render, screen } from "@testing-library/react";
import { InfiniteScroll } from "../infinite-scroll";

describe("InfiniteScroll", () => {
  const defaultProps = {
    loading: false,
    hasMore: true,
    onLoadMore: jest.fn(),
  };

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
});
