import { render, screen } from "@testing-library/react";
import { EmptyState } from "../empty-state";

describe("EmptyState", () => {
  const defaultProps = {
    title: "No items found",
    description: "Try adjusting your search filters.",
    cta: <button>Add Item</button>,
  };

  it("should render title", () => {
    render(<EmptyState {...defaultProps} />);
    expect(
      screen.getByRole("heading", { name: "No items found" }),
    ).toBeInTheDocument();
  });

  it("should render description", () => {
    render(<EmptyState {...defaultProps} />);
    expect(
      screen.getByText("Try adjusting your search filters."),
    ).toBeInTheDocument();
  });

  it("should render cta", () => {
    render(<EmptyState {...defaultProps} />);
    expect(
      screen.getByRole("button", { name: "Add Item" }),
    ).toBeInTheDocument();
  });

  it("should not render image wrapper when image is omitted", () => {
    const { container } = render(<EmptyState {...defaultProps} />);
    // The root div should only contain h3, p, and the cta element (3 children)
    expect(container.firstElementChild?.children.length).toBe(3);
  });

  it("should render image when provided", () => {
    render(<EmptyState {...defaultProps} image={<span>🎉</span>} />);
    expect(screen.getByText("🎉")).toBeInTheDocument();
  });
});
