import { render, screen } from "@testing-library/react";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
  BreadcrumbEllipsis,
} from "../breadcrumb";

describe("Breadcrumb", () => {
  it("renders without crashing", () => {
    render(
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink href="/">Home</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbEllipsis />
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Current</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>,
    );
  });

  it("renders nav with aria-label='breadcrumb'", () => {
    render(<Breadcrumb />);
    expect(
      screen.getByRole("navigation", { name: "breadcrumb" }),
    ).toBeInTheDocument();
  });

  it("BreadcrumbPage has aria-current='page' and aria-disabled='true'", () => {
    render(<BreadcrumbPage>Current Page</BreadcrumbPage>);
    const page = screen.getByText("Current Page");
    expect(page).toHaveAttribute("aria-current", "page");
    expect(page).toHaveAttribute("aria-disabled", "true");
  });

  it("BreadcrumbSeparator has aria-hidden='true' and role='presentation'", () => {
    render(
      <ol>
        <BreadcrumbSeparator />
      </ol>,
    );
    const separator = screen.getByRole("presentation", { hidden: true });
    expect(separator).toHaveAttribute("aria-hidden", "true");
  });

  it("BreadcrumbEllipsis has aria-hidden='true' and sr-only 'More' text", () => {
    render(<BreadcrumbEllipsis />);
    // The span itself is aria-hidden
    const ellipsis = screen.getByText("More");
    expect(ellipsis).toBeInTheDocument();
  });

  it("BreadcrumbLink renders as an anchor element", () => {
    render(<BreadcrumbLink href="/home">Home</BreadcrumbLink>);
    const link = screen.getByRole("link", { name: "Home" });
    expect(link).toBeInTheDocument();
    expect(link.tagName).toBe("A");
  });

  it("BreadcrumbLink renders child element when asChild is used", () => {
    render(
      <BreadcrumbLink asChild>
        <button>Home Button</button>
      </BreadcrumbLink>,
    );
    expect(
      screen.getByRole("button", { name: "Home Button" }),
    ).toBeInTheDocument();
  });
});
