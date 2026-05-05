import { render, screen } from "@testing-library/react";
import { Badge } from "../badge";

describe("Badge", () => {
  it("should render with default variant", () => {
    render(<Badge>Default</Badge>);
    expect(screen.getByText("Default")).toBeInTheDocument();
  });

  it("should render text children", () => {
    render(<Badge>Hello Badge</Badge>);
    expect(screen.getByText("Hello Badge")).toBeInTheDocument();
  });

  it.each([
    "default",
    "secondary",
    "destructive",
    "success",
    "warning",
    "info",
    "outline",
  ] as const)("should render with variant=%s without crashing", (variant) => {
    render(<Badge variant={variant}>Label</Badge>);
    expect(screen.getByText("Label")).toBeInTheDocument();
  });

  it("should accept enableShadow prop without crashing", () => {
    render(<Badge enableShadow>Shadow Badge</Badge>);
    expect(screen.getByText("Shadow Badge")).toBeInTheDocument();
  });

  it("should render as anchor element when asChild is used with <a>", () => {
    render(
      <Badge asChild>
        <a href="#">Link Badge</a>
      </Badge>,
    );
    expect(
      screen.getByRole("link", { name: "Link Badge" }),
    ).toBeInTheDocument();
  });

  it.each(["sm", "default", "lg"] as const)(
    "should accept size=%s without crashing",
    (size) => {
      render(<Badge size={size}>Sized</Badge>);
      expect(screen.getByText("Sized")).toBeInTheDocument();
    },
  );
});
