import { render } from "@testing-library/react";
import { Skeleton } from "../skeleton";

describe("Skeleton", () => {
  it("should render without crashing", () => {
    const { container } = render(<Skeleton />);
    expect(container.firstChild).toBeInTheDocument();
  });

  it("should render with additional className", () => {
    const { container } = render(<Skeleton className="h-10 w-full" />);
    expect(container.firstChild).toHaveClass("h-10", "w-full");
  });

  it("should render children when provided", () => {
    const { container } = render(
      <Skeleton>
        <span>content</span>
      </Skeleton>,
    );
    expect(container.querySelector("span")).toBeInTheDocument();
  });
});
