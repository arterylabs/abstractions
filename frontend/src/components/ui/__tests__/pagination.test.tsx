import { render, screen } from "@testing-library/react";
import {
  Pagination,
  PaginationContent,
  PaginationLink,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationEllipsis,
} from "../pagination";

describe("Pagination", () => {
  it("should render nav with role=navigation and aria-label='pagination'", () => {
    render(<Pagination />);
    expect(
      screen.getByRole("navigation", { name: "pagination" }),
    ).toBeInTheDocument();
  });
});

describe("PaginationPrevious", () => {
  it("should render with aria-label='Go to previous page'", () => {
    render(<PaginationPrevious />);
    expect(screen.getByLabelText("Go to previous page")).toBeInTheDocument();
  });
});

describe("PaginationNext", () => {
  it("should render with aria-label='Go to next page'", () => {
    render(<PaginationNext />);
    expect(screen.getByLabelText("Go to next page")).toBeInTheDocument();
  });
});

describe("PaginationLink", () => {
  it("should have aria-current='page' when isActive=true", () => {
    render(<PaginationLink isActive>1</PaginationLink>);
    expect(screen.getByText("1")).toHaveAttribute("aria-current", "page");
  });

  it("should not have aria-current when isActive is not set", () => {
    render(<PaginationLink>2</PaginationLink>);
    expect(screen.getByText("2")).not.toHaveAttribute("aria-current");
  });
});

describe("PaginationContent", () => {
  it("should render children", () => {
    render(
      <PaginationContent>
        <PaginationItem>
          <PaginationLink>3</PaginationLink>
        </PaginationItem>
      </PaginationContent>,
    );
    expect(screen.getByText("3")).toBeInTheDocument();
  });
});

describe("PaginationEllipsis", () => {
  it("should render without crashing", () => {
    render(<PaginationEllipsis />);
    expect(screen.getByText("More pages")).toBeInTheDocument();
  });
});
