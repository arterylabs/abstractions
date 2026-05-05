import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Input } from "../input";

describe("Input", () => {
  it("should render input element", () => {
    render(<Input id="test-input" value="" onChange={jest.fn()} />);
    expect(screen.getByRole("textbox")).toBeInTheDocument();
  });

  it("should render label when label prop is provided", () => {
    render(
      <Input id="test-input" value="" onChange={jest.fn()} label="Username" />,
    );
    expect(screen.getByText("Username")).toBeInTheDocument();
  });

  it("should call onChange with the typed string value", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<Input id="test-input" value="" onChange={onChange} />);
    await user.type(screen.getByRole("textbox"), "h");
    expect(onChange).toHaveBeenCalledWith("h");
  });

  it("should be disabled when disabled prop is provided", () => {
    render(<Input id="test-input" value="" onChange={jest.fn()} disabled />);
    expect(screen.getByRole("textbox")).toBeDisabled();
  });

  it("should show eye icon when toggleVisibility and type='password'", () => {
    render(
      <Input
        id="test-input"
        value=""
        onChange={jest.fn()}
        type="password"
        toggleVisibility
      />,
    );
    // Icon renders the type name as its text content; initially shows "visibility"
    expect(screen.getByText("visibility")).toBeInTheDocument();
  });

  it("should toggle input type between password and text when eye icon is clicked", async () => {
    const user = userEvent.setup();
    render(
      <Input
        id="test-input"
        value=""
        onChange={jest.fn()}
        type="password"
        toggleVisibility
      />,
    );
    // password inputs don't have role="textbox", query directly
    const passwordInput = document.querySelector("input") as HTMLInputElement;
    expect(passwordInput).toHaveAttribute("type", "password");

    await user.click(screen.getByText("visibility"));
    expect(passwordInput).toHaveAttribute("type", "text");
  });

  it("should render prepend content when provided", () => {
    render(
      <Input
        id="test-input"
        value=""
        onChange={jest.fn()}
        prepend={<span>$</span>}
      />,
    );
    expect(screen.getByText("$")).toBeInTheDocument();
  });

  it("should apply compact sizing when size is provided", () => {
    render(<Input id="test-input" value="" onChange={jest.fn()} size="sm" />);
    expect(screen.getByRole("textbox")).toHaveClass("h-8");
  });
});
