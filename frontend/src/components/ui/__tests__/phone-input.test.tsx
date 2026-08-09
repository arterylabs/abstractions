import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { PhoneInput } from "../phone-input";

describe("PhoneInput", () => {
  beforeEach(() => {
    // Make the IP-detection fetch reject immediately so the async state
    // update (setDialCode) is flushed before each test assertion.
    global.fetch = jest.fn().mockRejectedValue(new Error("offline"));
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it("renders without crashing", async () => {
    await act(async () => {
      render(<PhoneInput id="phone" onChange={jest.fn()} />);
    });
  });

  it("renders label when provided", async () => {
    await act(async () => {
      render(
        <PhoneInput id="phone" label="Phone Number" onChange={jest.fn()} />,
      );
    });
    expect(screen.getByText("Phone Number")).toBeInTheDocument();
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

  it("applies compact sizing when size is provided", async () => {
    await act(async () => {
      render(<PhoneInput id="phone" size="sm" onChange={jest.fn()} />);
    });

    expect(screen.getByRole("textbox").parentElement).toHaveClass("h-8");
  });

  it("emits E.164 when the local number changes", async () => {
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
  });

  it("re-emits with the new dial code when the country changes", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();

    await act(async () => {
      render(
        <PhoneInput id="phone" value="+911234567890" onChange={onChange} />,
      );
    });

    await user.selectOptions(screen.getByLabelText("Country dial code"), "+1");

    expect(onChange).toHaveBeenCalledWith("+11234567890");
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
