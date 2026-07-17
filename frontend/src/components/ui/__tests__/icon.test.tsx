import { render, screen } from "@testing-library/react";
import Icon from "../icon";

describe("Icon", () => {
  it("should render the type text as content", () => {
    render(<Icon type="home" />);
    expect(screen.getByText("home")).toBeInTheDocument();
  });

  it("should render with provided className", () => {
    const { container } = render(
      <Icon type="settings" className="text-red-500" />,
    );
    expect(container.querySelector("span")).toHaveClass("text-red-500");
  });

  it("should render without crashing when className is not provided", () => {
    render(<Icon type="close" />);
    expect(screen.getByText("close")).toBeInTheDocument();
  });
});
