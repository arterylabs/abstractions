import { render } from "@testing-library/react";
import { Progress } from "../progress";

describe("Progress", () => {
  it("should render without crashing", () => {
    const { container } = render(<Progress />);
    expect(container.firstChild).toBeInTheDocument();
  });

  it("should render with value prop", () => {
    const { container } = render(<Progress value={50} />);
    expect(
      container.querySelector("[data-slot='progress']"),
    ).toBeInTheDocument();
  });

  it("should render with value=0", () => {
    const { container } = render(<Progress value={0} />);
    expect(
      container.querySelector("[data-slot='progress']"),
    ).toBeInTheDocument();
  });

  it("should render with value=100", () => {
    const { container } = render(<Progress value={100} />);
    expect(
      container.querySelector("[data-slot='progress']"),
    ).toBeInTheDocument();
  });
});
