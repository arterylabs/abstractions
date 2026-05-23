import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import CopyButton from "../copy-button";

describe("CopyButton", () => {
  let writeTextMock: jest.Mock;

  beforeEach(() => {
    writeTextMock = jest.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText: writeTextMock },
      writable: true,
      configurable: true,
    });
  });

  it("renders without crashing", () => {
    render(<CopyButton text="hello" />);
  });

  it("renders a button with aria-label 'Copy code'", () => {
    render(<CopyButton text="hello" />);
    expect(
      screen.getByRole("button", { name: /copy code/i }),
    ).toBeInTheDocument();
  });

  it("calls clipboard.writeText with the provided text on click", async () => {
    const user = userEvent.setup();
    writeTextMock = jest.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText: writeTextMock },
      writable: true,
      configurable: true,
    });
    render(<CopyButton text="const x = 1;" />);
    await user.click(screen.getByRole("button", { name: /copy code/i }));
    expect(writeTextMock).toHaveBeenCalledWith("const x = 1;");
  });

  it("accepts a custom className", () => {
    render(<CopyButton text="hello" className="custom-class" />);
    expect(screen.getByRole("button", { name: /copy code/i })).toHaveClass(
      "custom-class",
    );
  });
});
