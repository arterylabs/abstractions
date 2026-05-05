import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { PageHeader } from "../page-header";

describe("PageHeader", () => {
  it("should render h1 with title", () => {
    render(<PageHeader title="Page Title" />);
    expect(
      screen.getByRole("heading", { level: 1, name: "Page Title" }),
    ).toBeInTheDocument();
  });

  it("should render description paragraph when provided", () => {
    render(<PageHeader title="Title" description="Some description" />);
    expect(screen.getByText("Some description")).toBeInTheDocument();
  });

  it("should not render back button by default", () => {
    render(<PageHeader title="Title" />);
    expect(screen.queryByRole("button")).not.toBeInTheDocument();
  });

  it("should render back button when showBackButton is true", () => {
    render(<PageHeader title="Title" showBackButton handleBack={jest.fn()} />);
    expect(screen.getByRole("button")).toBeInTheDocument();
  });

  it("should call handleBack when back button is clicked", async () => {
    const user = userEvent.setup();
    const handleBack = jest.fn();
    render(<PageHeader title="Title" showBackButton handleBack={handleBack} />);
    await user.click(screen.getByRole("button"));
    expect(handleBack).toHaveBeenCalledTimes(1);
  });

  it("should render children", () => {
    render(
      <PageHeader title="Title">
        <span>Child Content</span>
      </PageHeader>,
    );
    expect(screen.getByText("Child Content")).toBeInTheDocument();
  });
});
