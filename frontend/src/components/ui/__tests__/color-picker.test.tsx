import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ColorPicker } from "../color-picker";

// HexColorPicker doesn't play well with JSDOM; stub it out.
jest.mock("react-colorful", () => ({
  HexColorPicker: ({
    color,
    onChange,
  }: {
    color: string;
    onChange: (c: string) => void;
  }) => (
    <div data-testid="hex-color-picker" data-color={color}>
      <button onClick={() => onChange("#abcdef")}>pick</button>
    </div>
  ),
}));

describe("ColorPicker", () => {
  it("renders the swatch button with the current color", () => {
    render(<ColorPicker value="#ff0000" onChange={jest.fn()} />);
    expect(
      screen.getByRole("button", { name: /pick color/i }),
    ).toBeInTheDocument();
  });

  it("renders the hex text input with the current value", () => {
    render(<ColorPicker value="#ff0000" onChange={jest.fn()} />);
    expect(screen.getByDisplayValue("#ff0000")).toBeInTheDocument();
  });

  it("calls onChange with a normalised hex when a valid 6-digit value is typed", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<ColorPicker value="#000000" onChange={onChange} />);
    const input = screen.getByDisplayValue("#000000");
    await user.clear(input);
    await user.type(input, "#aabbcc");
    expect(onChange).toHaveBeenCalledWith("#aabbcc");
  });

  it("prefixes '#' and calls onChange when user omits the hash", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<ColorPicker value="#000000" onChange={onChange} />);
    const input = screen.getByDisplayValue("#000000");
    await user.clear(input);
    await user.type(input, "ff0000");
    expect(onChange).toHaveBeenCalledWith("#ff0000");
  });

  it("resets the input to the last valid value on blur when the text is invalid", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<ColorPicker value="#123456" onChange={onChange} />);
    const input = screen.getByDisplayValue("#123456");
    await user.clear(input);
    await user.type(input, "bad");
    fireEvent.blur(input);
    expect(screen.getByDisplayValue("#123456")).toBeInTheDocument();
  });

  it("calls onChange when the swatch picker fires a colour change", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    // Open the popover first
    render(<ColorPicker value="#000000" onChange={onChange} />);
    await user.click(screen.getByRole("button", { name: /pick color/i }));
    await user.click(screen.getByRole("button", { name: "pick" }));
    expect(onChange).toHaveBeenCalledWith("#abcdef");
  });

  it("syncs the text input when the parent changes the value prop", () => {
    const { rerender } = render(
      <ColorPicker value="#000000" onChange={jest.fn()} />,
    );
    rerender(<ColorPicker value="#ffffff" onChange={jest.fn()} />);
    expect(screen.getByDisplayValue("#ffffff")).toBeInTheDocument();
  });
});
