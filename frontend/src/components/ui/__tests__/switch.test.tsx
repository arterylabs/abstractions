import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Switch } from "../switch";

describe("Switch", () => {
  it("renders without crashing", () => {
    render(<Switch />);
  });

  it("renders a switch element", () => {
    render(<Switch />);
    expect(screen.getByRole("switch")).toBeInTheDocument();
  });

  it("renders in checked state when checked prop is true", () => {
    render(<Switch checked onCheckedChange={() => {}} />);
    expect(screen.getByRole("switch")).toHaveAttribute("aria-checked", "true");
  });

  it("calls onCheckedChange when toggled", async () => {
    const user = userEvent.setup();
    const onCheckedChange = jest.fn();
    render(<Switch onCheckedChange={onCheckedChange} />);
    await user.click(screen.getByRole("switch"));
    expect(onCheckedChange).toHaveBeenCalledTimes(1);
  });

  it("does not call onCheckedChange when disabled", async () => {
    const user = userEvent.setup();
    const onCheckedChange = jest.fn();
    render(<Switch disabled onCheckedChange={onCheckedChange} />);
    await user.click(screen.getByRole("switch"));
    expect(onCheckedChange).not.toHaveBeenCalled();
  });

  it("renders as disabled when disabled prop is provided", () => {
    render(<Switch disabled />);
    expect(screen.getByRole("switch")).toBeDisabled();
  });
});
