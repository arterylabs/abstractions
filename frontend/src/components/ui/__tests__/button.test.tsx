import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "../button";

describe("Button", () => {
  it("should render a button element by default", () => {
    render(<Button>Click me</Button>);
    expect(
      screen.getByRole("button", { name: /Click me/i }),
    ).toBeInTheDocument();
  });

  it("should render children text", () => {
    render(<Button>Children Text</Button>);
    expect(screen.getByText("Children Text")).toBeInTheDocument();
  });

  it("should call onClick when clicked", async () => {
    const user = userEvent.setup();
    const onClick = jest.fn();
    render(<Button onClick={onClick}>Click me</Button>);
    await user.click(screen.getByRole("button", { name: /Click me/i }));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("should not call onClick when disabled", async () => {
    const user = userEvent.setup();
    const onClick = jest.fn();
    render(
      <Button onClick={onClick} disabled>
        Disabled
      </Button>,
    );
    await user.click(screen.getByRole("button", { name: /Disabled/i }));
    expect(onClick).not.toHaveBeenCalled();
  });

  it("should show Spinner (aria-hidden svg) when loading is true", () => {
    const { container } = render(<Button loading>Loading</Button>);
    expect(
      container.querySelector('svg[aria-hidden="true"]'),
    ).toBeInTheDocument();
  });

  it("should render as anchor element when asChild is used with <a>", () => {
    render(
      <Button asChild>
        <a href="#">Link Button</a>
      </Button>,
    );
    expect(
      screen.getByRole("link", { name: "Link Button" }),
    ).toBeInTheDocument();
  });

  it.each([
    "default",
    "destructive",
    "outline",
    "secondary",
    "ghost",
    "link",
  ] as const)("should render variant=%s without crashing", (variant) => {
    render(<Button variant={variant}>Variant</Button>);
    expect(
      screen.getByRole("button", { name: /Variant/i }),
    ).toBeInTheDocument();
  });
});
