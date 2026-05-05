import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Label } from "../label";

describe("Label", () => {
  it("renders without crashing", () => {
    render(<Label>My Label</Label>);
  });

  it("renders text content", () => {
    render(<Label>Full Name</Label>);
    expect(screen.getByText("Full Name")).toBeInTheDocument();
  });

  it("passes htmlFor to the underlying element", () => {
    render(<Label htmlFor="email-input">Email</Label>);
    const label = screen.getByText("Email");
    expect(label).toHaveAttribute("for", "email-input");
  });

  it("clicking label focuses associated input", async () => {
    const user = userEvent.setup();
    render(
      <>
        <Label htmlFor="name-field">Name</Label>
        <input id="name-field" />
      </>,
    );
    await user.click(screen.getByText("Name"));
    expect(screen.getByRole("textbox")).toHaveFocus();
  });
});
