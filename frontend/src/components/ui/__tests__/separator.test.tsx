import { render, screen } from "@testing-library/react";
import { Separator } from "../separator";

describe("Separator", () => {
  it("renders without crashing", () => {
    render(<Separator />);
  });

  it("renders as horizontal by default (data-orientation=horizontal)", () => {
    render(<Separator />);
    expect(screen.getByRole("none", { hidden: true })).toHaveAttribute(
      "data-orientation",
      "horizontal",
    );
  });

  it("renders as vertical when orientation='vertical'", () => {
    render(<Separator orientation="vertical" />);
    expect(screen.getByRole("none", { hidden: true })).toHaveAttribute(
      "data-orientation",
      "vertical",
    );
  });

  it("renders with role='separator' when decorative=false", () => {
    render(<Separator decorative={false} />);
    expect(screen.getByRole("separator")).toBeInTheDocument();
  });
});
