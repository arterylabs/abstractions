import { render, screen } from "@testing-library/react";
import { UserAvatar } from "../user-avatar";

describe("UserAvatar", () => {
  it("renders the initials fallback derived from the name", () => {
    render(<UserAvatar name="Jane Doe" />);
    expect(screen.getByText("JD")).toBeInTheDocument();
  });

  it("renders only the first initial for a single-word name", () => {
    render(<UserAvatar name="Alice" />);
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("renders without crashing when imageUrl is provided", () => {
    // In JSDOM images never load, so Radix Avatar.Image won't mount an <img>
    // element; the fallback stays visible. We just verify the component
    // renders without throwing.
    render(
      <UserAvatar name="Jane Doe" imageUrl="https://example.com/avatar.jpg" />,
    );
    // The avatar wrapper is always present
    expect(screen.getByText("JD")).toBeInTheDocument();
  });

  it("does not render an img when imageUrl is null", () => {
    render(<UserAvatar name="Jane Doe" imageUrl={null} />);
    expect(screen.queryByRole("img")).not.toBeInTheDocument();
  });

  it("applies a custom className to the avatar wrapper", () => {
    const { container } = render(
      <UserAvatar name="Jane Doe" className="h-12 w-12" />,
    );
    expect(container.firstChild).toHaveClass("h-12", "w-12");
  });
});
