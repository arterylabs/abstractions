import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button3D } from "../button-3d";

describe("Button3D", () => {
  it("renders children text", () => {
    render(<Button3D>Click me</Button3D>);
    expect(
      screen.getByRole("button", { name: /click me/i }),
    ).toBeInTheDocument();
  });

  it("applies the default (primary) variant class", () => {
    render(<Button3D>Default</Button3D>);
    expect(screen.getByRole("button", { name: /default/i })).toHaveClass(
      "bg-primary",
    );
  });

  it("applies the destructive variant class", () => {
    render(<Button3D variant="destructive">Delete</Button3D>);
    expect(screen.getByRole("button", { name: /delete/i })).toHaveClass(
      "bg-destructive",
    );
  });

  it("calls onClick when clicked", async () => {
    const user = userEvent.setup();
    const onClick = jest.fn();
    render(<Button3D onClick={onClick}>Press</Button3D>);
    await user.click(screen.getByRole("button", { name: /press/i }));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("shows Spinner and pointer-events-none when loading is true", () => {
    const { container } = render(<Button3D loading>Save</Button3D>);
    expect(
      container.querySelector('svg[aria-hidden="true"]'),
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /save/i })).toHaveClass(
      "pointer-events-none",
    );
  });

  it("renders as a Slot (anchor) when asChild is used", () => {
    render(
      <Button3D asChild>
        <a href="#">Link</a>
      </Button3D>,
    );
    expect(screen.getByRole("link", { name: "Link" })).toBeInTheDocument();
  });

  it.each([
    "default",
    "destructive",
    "outline",
    "secondary",
    "ghost",
    "link",
  ] as const)("renders variant=%s without crashing", (variant) => {
    render(<Button3D variant={variant}>Variant</Button3D>);
    expect(
      screen.getByRole("button", { name: /variant/i }),
    ).toBeInTheDocument();
  });
});
