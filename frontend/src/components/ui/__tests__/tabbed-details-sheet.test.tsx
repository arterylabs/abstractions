import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TabbedDetailsSheet } from "../tabbed-details-sheet";

const tabs = [
  { value: "overview", label: "Overview", content: <p>Overview content</p> },
  { value: "details", label: "Details", content: <p>Details content</p> },
  { value: "history", label: "History", content: <p>History content</p> },
];

describe("TabbedDetailsSheet", () => {
  it("renders the sheet title when open", () => {
    render(
      <TabbedDetailsSheet
        isOpen
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
      />,
    );
    expect(screen.getByText("My Sheet")).toBeInTheDocument();
  });

  it("renders all tab labels", () => {
    render(
      <TabbedDetailsSheet
        isOpen
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
      />,
    );
    expect(screen.getByRole("tab", { name: "Overview" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Details" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "History" })).toBeInTheDocument();
  });

  it("shows the first tab's content by default", () => {
    render(
      <TabbedDetailsSheet
        isOpen
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
      />,
    );
    expect(screen.getByText("Overview content")).toBeInTheDocument();
  });

  it("shows the defaultTab content when defaultTab is specified", () => {
    render(
      <TabbedDetailsSheet
        isOpen
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
        defaultTab="details"
      />,
    );
    expect(
      screen.getByRole("tab", { name: "Details", selected: true }),
    ).toBeInTheDocument();
  });

  it("switches tab content when a different tab is clicked", async () => {
    const user = userEvent.setup();
    render(
      <TabbedDetailsSheet
        isOpen
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
      />,
    );
    await user.click(screen.getByRole("tab", { name: "Details" }));
    expect(screen.getByText("Details content")).toBeInTheDocument();
  });

  it("does not render content when closed", () => {
    render(
      <TabbedDetailsSheet
        isOpen={false}
        onClose={jest.fn()}
        title="My Sheet"
        tabs={tabs}
      />,
    );
    expect(screen.queryByText("My Sheet")).not.toBeInTheDocument();
  });
});
