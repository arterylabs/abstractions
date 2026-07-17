import { render } from "@testing-library/react";
import Spinner from "../spinner";

describe("Spinner", () => {
  it("should render an SVG element", () => {
    const { container } = render(<Spinner />);
    expect(container.querySelector("svg")).toBeInTheDocument();
  });

  it("should have aria-hidden='true'", () => {
    const { container } = render(<Spinner />);
    expect(container.querySelector("svg")).toHaveAttribute(
      "aria-hidden",
      "true",
    );
  });

  it("should accept fill='primary' without crashing", () => {
    const { container } = render(<Spinner fill="primary" />);
    expect(container.querySelector("svg")).toBeInTheDocument();
  });

  it("should accept fill='secondary' without crashing", () => {
    const { container } = render(<Spinner fill="secondary" />);
    expect(container.querySelector("svg")).toBeInTheDocument();
  });
});
