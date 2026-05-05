import { render, screen } from "@testing-library/react";
import Logo from "../logo";

describe("Logo", () => {
  it("should render img with alt='Harmony Logo'", () => {
    render(<Logo />);
    expect(screen.getByAltText("Harmony Logo")).toBeInTheDocument();
  });

  it("should render icon src when icon=true", () => {
    render(<Logo icon={true} />);
    expect(screen.getByAltText("Harmony Logo")).toHaveAttribute(
      "src",
      "/assets/icon.svg",
    );
  });

  it("should render full logo src when icon=false", () => {
    render(<Logo icon={false} />);
    expect(screen.getByAltText("Harmony Logo")).toHaveAttribute(
      "src",
      "/assets/logo.svg",
    );
  });
});
