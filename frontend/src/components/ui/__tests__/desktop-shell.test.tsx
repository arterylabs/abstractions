import { render, screen } from "@testing-library/react";
import DesktopShell from "../desktop-shell";

describe("DesktopShell", () => {
  it("renders children inside the shell", () => {
    render(
      <DesktopShell>
        <p>Page content</p>
      </DesktopShell>,
    );
    expect(screen.getByText("Page content")).toBeInTheDocument();
  });

  it("displays the provided URL in the address bar", () => {
    render(
      <DesktopShell url="example.com/page">
        <span>content</span>
      </DesktopShell>,
    );
    // URL appears in both the tab strip and the address bar
    expect(
      screen.getAllByText("example.com/page").length,
    ).toBeGreaterThanOrEqual(1);
  });

  it("falls back to 'harmony.app' when no URL is given", () => {
    render(
      <DesktopShell>
        <span>content</span>
      </DesktopShell>,
    );
    expect(screen.getAllByText("harmony.app").length).toBeGreaterThanOrEqual(1);
  });
});
