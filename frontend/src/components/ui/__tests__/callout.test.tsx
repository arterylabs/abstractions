import { render, screen } from "@testing-library/react";
import Callout from "../callout";

describe("Callout", () => {
  it("renders without crashing", () => {
    render(<Callout>Some content</Callout>);
  });

  it("renders children content", () => {
    render(<Callout>This is a note</Callout>);
    expect(screen.getByText("This is a note")).toBeInTheDocument();
  });

  it("renders default type 'note' with default title", () => {
    render(<Callout>Content</Callout>);
    expect(screen.getByText("Note")).toBeInTheDocument();
  });

  it("renders custom title when provided", () => {
    render(<Callout title="Custom Title">Content</Callout>);
    expect(screen.getByText("Custom Title")).toBeInTheDocument();
  });

  it.each(["note", "tip", "warning", "danger"] as const)(
    "renders type=%s without crashing",
    (type) => {
      render(<Callout type={type}>Content</Callout>);
      expect(screen.getByText("Content")).toBeInTheDocument();
    },
  );

  it("renders 'Tip' as default title for type=tip", () => {
    render(<Callout type="tip">Content</Callout>);
    expect(screen.getByText("Tip")).toBeInTheDocument();
  });

  it("renders 'Warning' as default title for type=warning", () => {
    render(<Callout type="warning">Content</Callout>);
    expect(screen.getByText("Warning")).toBeInTheDocument();
  });

  it("renders 'Danger' as default title for type=danger", () => {
    render(<Callout type="danger">Content</Callout>);
    expect(screen.getByText("Danger")).toBeInTheDocument();
  });
});
