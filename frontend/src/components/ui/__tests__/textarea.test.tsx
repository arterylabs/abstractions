import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Textarea } from "../textarea";

describe("Textarea", () => {
  it("renders without crashing", () => {
    render(<Textarea />);
  });

  it("renders a textarea element", () => {
    render(<Textarea />);
    expect(screen.getByRole("textbox")).toBeInTheDocument();
  });

  it("renders with placeholder text", () => {
    render(<Textarea placeholder="Enter notes here" />);
    expect(screen.getByPlaceholderText("Enter notes here")).toBeInTheDocument();
  });

  it("renders controlled value", () => {
    render(<Textarea value="Hello world" onChange={() => {}} />);
    expect(screen.getByRole("textbox")).toHaveValue("Hello world");
  });

  it("calls onChange when typing", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<Textarea onChange={onChange} />);
    await user.type(screen.getByRole("textbox"), "test");
    expect(onChange).toHaveBeenCalled();
  });

  it("is disabled when disabled prop is provided", () => {
    render(<Textarea disabled />);
    expect(screen.getByRole("textbox")).toBeDisabled();
  });

  it("passes aria-invalid attribute to the element", () => {
    render(<Textarea aria-invalid="true" />);
    expect(screen.getByRole("textbox")).toHaveAttribute("aria-invalid", "true");
  });
});
