import { useState } from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { DeleteRepositoryDialog } from "@/app/(dashboard)/repositories/[id]/settings/DeleteRepositoryDialog";
import {
  isoWeekdayToName,
  localHourToUtc,
  utcHourToLocal,
} from "@/lib/timezone";

describe("repository settings", () => {
  it("converts UTC 0 to a local time string", () => {
    const now = new Date();
    const expected = new Date(
      Date.UTC(now.getFullYear(), now.getMonth(), now.getDate(), 0),
    ).toLocaleTimeString([], {
      hour: "numeric",
      minute: "2-digit",
    });

    expect(utcHourToLocal(0)).toBe(expected);
  });

  it("round-trips localHourToUtc correctly", () => {
    const now = new Date();
    const utcHour = 9;
    const localHour = new Date(
      Date.UTC(now.getFullYear(), now.getMonth(), now.getDate(), utcHour),
    ).getHours();

    expect(localHourToUtc(localHour)).toBe(utcHour);
  });

  it("returns Monday for ISO weekday 1", () => {
    expect(isoWeekdayToName(1)).toBe("Monday");
  });

  it("blocks delete until the repository name matches", async () => {
    const user = userEvent.setup();
    const onConfirm = jest.fn();

    function TestDialog() {
      const [confirmationValue, setConfirmationValue] = useState("");

      return (
        <DeleteRepositoryDialog
          open
          onOpenChange={jest.fn()}
          repositoryFullName="acme/repo"
          confirmationValue={confirmationValue}
          onConfirmationChange={setConfirmationValue}
          onConfirm={onConfirm}
          isDeleting={false}
        />
      );
    }

    render(<TestDialog />);

    const input = screen.getByLabelText("Repository name");
    const confirmButton = screen.getByRole("button", {
      name: "Delete Repository",
    });

    expect(confirmButton).toBeDisabled();

    await user.type(input, "wrong/name");
    expect(confirmButton).toBeDisabled();

    await user.clear(input);
    await user.type(input, "acme/repo");

    expect(confirmButton).toBeEnabled();
  });
});
