import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
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

  it("shows the formatted date when a value is provided", () => {
    render(
      <DateTimePicker
        value={new Date("2024-01-15T10:30:00")}
        onChange={jest.fn()}
      />,
    );
    // Button label shows the formatted date (e.g. "January 15th, 2024, 10:30 AM")
    expect(screen.getByRole("button")).toHaveTextContent("2024");
  });

  it("opens the popover and shows the calendar when the trigger is clicked", async () => {
    const user = userEvent.setup();
    render(<DateTimePicker onChange={jest.fn()} />);

    await user.click(screen.getByRole("button"));

    // Calendar renders a grid (month view table)
    expect(screen.getByRole("grid")).toBeInTheDocument();
  });

  it("calls onChange with an updated time when the time input changes", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    const value = new Date("2024-01-15T10:30:00");

    render(<DateTimePicker value={value} onChange={onChange} id="picker" />);

    // Open the popover
    await user.click(screen.getByRole("button"));

    // The time input should display the current time "10:30"
    const timeInput = screen.getByDisplayValue("10:30");

    // fireEvent.change is used here because userEvent does not reliably
    // simulate native <input type="time"> value changes in jsdom
    fireEvent.change(timeInput, { target: { value: "14:45" } });

    expect(onChange).toHaveBeenCalled();
    const result = onChange.mock.calls[0][0] as Date;
    expect(result).toBeInstanceOf(Date);
    expect(result.getHours()).toBe(14);
    expect(result.getMinutes()).toBe(45);
  });

  it("closes the popover when the Done button is clicked", async () => {
    const user = userEvent.setup();
    render(<DateTimePicker onChange={jest.fn()} />);

    await user.click(screen.getByRole("button")); // open
    expect(screen.getByRole("grid")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Done" }));
    expect(screen.queryByRole("grid")).not.toBeInTheDocument();
  });
});
