import { render, screen } from "@testing-library/react";
import { EmptyState } from "../empty-state";

describe("EmptyState", () => {
  const defaultProps = {
    icon: "inbox",
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

  it("should render icon", () => {
    render(<EmptyState {...defaultProps} />);
    expect(screen.getByText("inbox")).toBeInTheDocument();
  });

  it("should always render icon wrapper as first child", () => {
    const { container } = render(<EmptyState {...defaultProps} />);
    // The root div should contain icon wrapper, h3, p, and cta element (4 children)
    expect(container.firstElementChild?.children.length).toBe(4);
  });
});
