import { render, screen } from "@testing-library/react";
import Tooltip from "../tooltip";

describe("Tooltip", () => {
  it("should render children", () => {
    render(
      <Tooltip value="Tooltip text">
        <button>Hover me</button>
      </Tooltip>,
    );
    expect(
      screen.getByRole("button", { name: "Hover me" }),
    ).toBeInTheDocument();
  });

  it("should not throw when value is a string", () => {
    expect(() => {
      render(
        <Tooltip value="Simple tooltip text">
          <span>Target</span>
        </Tooltip>,
      );
    }).not.toThrow();
  });

  it("should not throw when value is a ReactNode", () => {
    expect(() => {
      render(
        <Tooltip value={<div>Rich tooltip content</div>}>
          <span>Target</span>
        </Tooltip>,
      );
    }).not.toThrow();
  });
});
