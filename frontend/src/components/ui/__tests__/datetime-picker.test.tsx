import { render } from "@testing-library/react";
import { DateTimePicker } from "../datetime-picker";

describe("DateTimePicker", () => {
  it("renders without crashing", () => {
    render(<DateTimePicker onChange={jest.fn()} />);
  });

  it("renders with a pre-selected date without crashing", () => {
    render(
      <DateTimePicker value={new Date("2024-01-15")} onChange={jest.fn()} />,
    );
  });
});
