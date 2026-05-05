import { render, screen } from "@testing-library/react";
import Logo from "../logo";

describe("Logo", () => {
  it("should render img with alt='Abstractions Logo'", () => {
    render(<Logo />);
    expect(screen.getByAltText("Abstractions Logo")).toBeInTheDocument();
  });

  it("should render icon src when icon=true", () => {
    render(<Logo icon={true} />);
    expect(screen.getByAltText("Abstractions Logo")).toHaveAttribute(
      "src",
      "/assets/icon.svg",
    );
  });

  it("should render full logo src when icon=false", () => {
    render(<Logo icon={false} />);
    expect(screen.getByAltText("Abstractions Logo")).toHaveAttribute(
      "src",
      "/assets/logo.svg",
    );
  });
});
