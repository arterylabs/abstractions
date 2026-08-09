import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { PhoneInput } from "../phone-input";

describe("PhoneInput", () => {
  beforeEach(() => {
    // Reject IP geo lookup so async setDialCode settles before assertions.
    global.fetch = jest.fn().mockRejectedValue(new Error("offline"));
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it("hydrates from an existing full phone value", async () => {
    await act(async () => {
      render(
        <PhoneInput id="phone" value="+919876543210" onChange={jest.fn()} />,
      );
    });

    expect(screen.getByDisplayValue("9876543210")).toBeInTheDocument();
    expect(screen.getByLabelText("Country dial code")).toHaveValue("+91");
  });

  it("emits E.164 when the local number or dial code changes", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();

    await act(async () => {
      render(
        <PhoneInput id="phone" value="+911234567890" onChange={onChange} />,
      );
    });

    const input = screen.getByRole("textbox");
    await user.clear(input);
    await user.type(input, "9876543210");
    expect(onChange).toHaveBeenLastCalledWith("+919876543210");

    await user.selectOptions(screen.getByLabelText("Country dial code"), "+1");
    expect(onChange).toHaveBeenCalledWith("+19876543210");
  });

  it("emits empty string when the local number is cleared", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();

    await act(async () => {
      render(
        <PhoneInput id="phone" value="+919876543210" onChange={onChange} />,
      );
    });

    await user.clear(screen.getByRole("textbox"));
    expect(onChange).toHaveBeenCalledWith("");
  });

  it("shows a validation error after blur when the number is invalid", async () => {
    const user = userEvent.setup();

    await act(async () => {
      render(
        <PhoneInput id="phone" value="+911234567890" onChange={jest.fn()} />,
      );
    });

    const input = screen.getByRole("textbox");
    await user.clear(input);
    await user.type(input, "12");
    await user.tab();

    expect(
      screen.getByText("Enter a valid phone number (digits only)."),
    ).toBeInTheDocument();
  });
});
