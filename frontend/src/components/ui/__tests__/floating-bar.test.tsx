import { render, screen } from "@testing-library/react";
import { FloatingBar } from "../floating-bar";

describe("FloatingBar", () => {
  const defaultProps = {
    content: <span>Unsaved changes</span>,
    saveButton: <button>Save</button>,
    cancelButton: <button>Cancel</button>,
  };

  it("should render content", () => {
    render(<FloatingBar {...defaultProps} />);
    expect(screen.getByText("Unsaved changes")).toBeInTheDocument();
  });

  it("should render saveButton", () => {
    render(<FloatingBar {...defaultProps} />);
    expect(screen.getByRole("button", { name: "Save" })).toBeInTheDocument();
  });

  it("should render cancelButton", () => {
    render(<FloatingBar {...defaultProps} />);
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument();
  });

  it("should render without crashing when visible=false", () => {
    render(<FloatingBar {...defaultProps} visible={false} />);
    expect(screen.getByText("Unsaved changes")).toBeInTheDocument();
  });

  it("should render without crashing when visible=true", () => {
    render(<FloatingBar {...defaultProps} visible={true} />);
    expect(screen.getByText("Unsaved changes")).toBeInTheDocument();
  });
});
