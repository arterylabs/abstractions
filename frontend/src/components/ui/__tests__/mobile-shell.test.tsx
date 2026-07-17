import { render, screen } from "@testing-library/react";
import MobileShell from "../mobile-shell";

describe("MobileShell", () => {
  it("renders children inside the shell", () => {
    render(
      <MobileShell>
        <p>App content</p>
      </MobileShell>,
    );
    expect(screen.getByText("App content")).toBeInTheDocument();
  });

  it("renders the status bar time indicator", () => {
    render(
      <MobileShell>
        <span>content</span>
      </MobileShell>,
    );
    expect(screen.getByText("9:41")).toBeInTheDocument();
  });

  it("renders without crashing when no children are provided", () => {
    const { container } = render(<MobileShell>{null}</MobileShell>);
    expect(container.firstChild).toBeInTheDocument();
  });
});
