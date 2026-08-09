import { render, screen } from "@testing-library/react";
import { SectionMark } from "../section-mark";

describe("SectionMark", () => {
  it("should render the label", () => {
    render(<SectionMark label="Profile" />);
    expect(screen.getByText("Profile")).toBeInTheDocument();
  });

  it("should render the index when provided", () => {
    render(<SectionMark label="Providers" index="01" />);
    expect(screen.getByText("§ 01")).toBeInTheDocument();
    expect(screen.getByText("Providers")).toBeInTheDocument();
  });

  it("should not render an index when omitted", () => {
    render(<SectionMark label="General" />);
    expect(screen.queryByText(/^§/)).not.toBeInTheDocument();
  });

  it("should render the action when provided", () => {
    render(
      <SectionMark
        label="Invitations"
        action={<button type="button">Invite</button>}
      />,
    );
    expect(screen.getByRole("button", { name: "Invite" })).toBeInTheDocument();
  });

  it("should not render an action container when action is omitted", () => {
    const { container } = render(<SectionMark label="Members" />);
    // Root has: optional index, label, hairline rule — no action wrapper
    expect(container.querySelectorAll("button")).toHaveLength(0);
  });

  it("should apply a custom className to the root", () => {
    const { container } = render(
      <SectionMark label="Danger zone" className="mt-4" />,
    );
    expect(container.firstElementChild).toHaveClass("mt-4");
  });
});
