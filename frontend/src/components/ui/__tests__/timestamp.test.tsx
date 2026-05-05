import { render, screen } from "@testing-library/react";
import { format } from "date-fns";
import Timestamp from "../timestamp";

describe("Timestamp", () => {
  it("should render a span element for a valid date", () => {
    const { container } = render(
      <Timestamp timestamp="2023-06-15T12:00:00Z" />,
    );
    expect(container.querySelector("span")).toBeInTheDocument();
  });

  it("should render 'Invalid date' for an invalid date string", () => {
    render(<Timestamp timestamp="not-a-valid-date" />);
    expect(screen.getByText("Invalid date")).toBeInTheDocument();
  });

  it("should render a formatted date for a valid ISO date with format 'PP'", () => {
    const date = new Date("2023-06-15T12:00:00Z");
    const expected = format(date, "PP");
    render(<Timestamp timestamp={date} format="PP" />);
    expect(screen.getByText(expected)).toBeInTheDocument();
  });

  it("should render a non-empty string for the default time-ago format", () => {
    render(<Timestamp timestamp={new Date().toISOString()} />);
    const span = document.querySelector("span");
    expect(span?.textContent?.trim().length).toBeGreaterThan(0);
  });
});
