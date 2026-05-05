import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  Collapsible,
  CollapsibleTrigger,
  CollapsibleContent,
} from "../collapsible";

describe("Collapsible", () => {
  it("renders without crashing", () => {
    render(
      <Collapsible>
        <CollapsibleTrigger>Toggle</CollapsibleTrigger>
        <CollapsibleContent>Content</CollapsibleContent>
      </Collapsible>,
    );
  });

  it("CollapsibleContent is hidden by default (open=false)", () => {
    render(
      <Collapsible>
        <CollapsibleTrigger>Toggle</CollapsibleTrigger>
        <CollapsibleContent>Hidden Content</CollapsibleContent>
      </Collapsible>,
    );
    // Radix sets data-state=closed on the content when closed
    const content = screen.queryByText("Hidden Content");
    if (content) {
      expect(content.closest("[data-state]")).toHaveAttribute(
        "data-state",
        "closed",
      );
    }
  });

  it("CollapsibleContent is visible when open=true", () => {
    render(
      <Collapsible open>
        <CollapsibleTrigger>Toggle</CollapsibleTrigger>
        <CollapsibleContent>Visible Content</CollapsibleContent>
      </Collapsible>,
    );
    expect(screen.getByText("Visible Content")).toBeInTheDocument();
    expect(
      screen.getByText("Visible Content").closest("[data-state]"),
    ).toHaveAttribute("data-state", "open");
  });

  it("clicking CollapsibleTrigger opens the content", async () => {
    const user = userEvent.setup();
    render(
      <Collapsible>
        <CollapsibleTrigger>Toggle</CollapsibleTrigger>
        <CollapsibleContent>Toggled Content</CollapsibleContent>
      </Collapsible>,
    );
    await user.click(screen.getByText("Toggle"));
    expect(
      screen.getByText("Toggled Content").closest("[data-state]"),
    ).toHaveAttribute("data-state", "open");
  });
});
