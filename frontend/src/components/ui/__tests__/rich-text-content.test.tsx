import { render, screen } from "@testing-library/react";
import { RichTextContent } from "../rich-text-content";

describe("RichTextContent", () => {
  it("renders nothing when html is null", () => {
    const { container } = render(<RichTextContent html={null} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders nothing when html is undefined", () => {
    const { container } = render(<RichTextContent />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders html content", () => {
    render(<RichTextContent html="<p>Hello world</p>" />);
    expect(screen.getByText("Hello world")).toBeInTheDocument();
  });

  it("renders with custom className", () => {
    const { container } = render(
      <RichTextContent html="<p>Text</p>" className="custom-class" />,
    );
    expect(container.firstChild).toHaveClass("custom-class");
  });
});
