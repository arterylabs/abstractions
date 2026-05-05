import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RadioGroup, RadioGroupItem } from "../radio-group";

describe("RadioGroup", () => {
  it("renders without crashing", () => {
    render(
      <RadioGroup defaultValue="option-1">
        <RadioGroupItem value="option-1" />
        <RadioGroupItem value="option-2" />
      </RadioGroup>,
    );
  });

  it("renders RadioGroupItem elements", () => {
    render(
      <RadioGroup>
        <RadioGroupItem value="a" aria-label="Option A" />
        <RadioGroupItem value="b" aria-label="Option B" />
      </RadioGroup>,
    );
    expect(screen.getAllByRole("radio")).toHaveLength(2);
  });

  it("fires onValueChange when an item is selected", async () => {
    const user = userEvent.setup();
    const onValueChange = jest.fn();
    render(
      <RadioGroup onValueChange={onValueChange}>
        <RadioGroupItem value="yes" aria-label="Yes" />
        <RadioGroupItem value="no" aria-label="No" />
      </RadioGroup>,
    );
    await user.click(screen.getByRole("radio", { name: "Yes" }));
    expect(onValueChange).toHaveBeenCalledWith("yes");
  });

  it("disabled RadioGroupItem cannot be interacted with", async () => {
    const user = userEvent.setup();
    const onValueChange = jest.fn();
    render(
      <RadioGroup onValueChange={onValueChange}>
        <RadioGroupItem value="disabled-opt" aria-label="Disabled" disabled />
      </RadioGroup>,
    );
    await user.click(screen.getByRole("radio", { name: "Disabled" }));
    expect(onValueChange).not.toHaveBeenCalled();
  });

  it("renders with a pre-selected default value", () => {
    render(
      <RadioGroup defaultValue="option-2">
        <RadioGroupItem value="option-1" aria-label="One" />
        <RadioGroupItem value="option-2" aria-label="Two" />
      </RadioGroup>,
    );
    expect(screen.getByRole("radio", { name: "Two" })).toBeChecked();
    expect(screen.getByRole("radio", { name: "One" })).not.toBeChecked();
  });
});
