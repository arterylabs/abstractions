import { act, render, screen } from "@testing-library/react";
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
});
