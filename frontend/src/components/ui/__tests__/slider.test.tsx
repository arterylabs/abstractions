import { render } from "@testing-library/react";
import { Slider } from "../slider";

describe("Slider", () => {
  it("renders without crashing", () => {
    render(<Slider defaultValue={[50]} />);
  });

  it("renders two thumbs by default when no value or defaultValue is provided (uses [min, max])", () => {
    const { container } = render(<Slider min={0} max={100} />);
    const thumbs = container.querySelectorAll("[data-slot='slider-thumb']");
    expect(thumbs).toHaveLength(2);
  });

  it("renders one thumb when value is a single-element array", () => {
    const { container } = render(
      <Slider value={[40]} onValueChange={() => {}} />,
    );
    const thumbs = container.querySelectorAll("[data-slot='slider-thumb']");
    expect(thumbs).toHaveLength(1);
  });

  it("renders two thumbs when value is a two-element array (range slider)", () => {
    const { container } = render(
      <Slider value={[20, 80]} onValueChange={() => {}} />,
    );
    const thumbs = container.querySelectorAll("[data-slot='slider-thumb']");
    expect(thumbs).toHaveLength(2);
  });

  it("renders one thumb when defaultValue is a single-element array", () => {
    const { container } = render(<Slider defaultValue={[50]} />);
    const thumbs = container.querySelectorAll("[data-slot='slider-thumb']");
    expect(thumbs).toHaveLength(1);
  });
});
