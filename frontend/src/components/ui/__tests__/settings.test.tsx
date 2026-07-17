import { useState } from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { DeleteRepositoryDialog } from "@/app/(dashboard)/repositories/[id]/settings/DeleteRepositoryDialog";
import { isoWeekdayToName, utcDayHourToLocal } from "@/lib/timezone";

describe("repository settings", () => {
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

describe("utcDayHourToLocal", () => {
  // All tests use DST-free timezones so expectations are stable year-round.

  it("applies a positive offset without crossing midnight", () => {
    // Asia/Tokyo = UTC+9 (no DST): UTC Wed 10:00 → local Wed 19:00
    const result = utcDayHourToLocal(3, 10, "Asia/Tokyo");
    expect(result.localHour).toBe(19);
    expect(result.localDay).toBe(3);
  });

  it("applies a negative offset without crossing midnight", () => {
    // Pacific/Honolulu = UTC-10 (no DST): UTC Wed 15:00 → local Wed 05:00
    const result = utcDayHourToLocal(3, 15, "Pacific/Honolulu");
    expect(result.localHour).toBe(5);
    expect(result.localDay).toBe(3);
  });

  it("rolls localHour forward and increments day when crossing midnight ahead", () => {
    // Asia/Tokyo = UTC+9: UTC Wed 20:00 → local Thu 05:00
    const result = utcDayHourToLocal(3, 20, "Asia/Tokyo");
    expect(result.localHour).toBe(5);
    expect(result.localDay).toBe(4); // Thursday
  });

  it("rolls localHour back and decrements day when crossing midnight behind", () => {
    // Pacific/Honolulu = UTC-10: UTC Mon 05:00 → local Sun 19:00
    const result = utcDayHourToLocal(1, 5, "Pacific/Honolulu");
    expect(result.localHour).toBe(19);
    expect(result.localDay).toBe(7); // Sunday wraps from Monday
  });

  it("wraps day correctly at the Sunday→Monday boundary", () => {
    // Asia/Tokyo = UTC+9: UTC Sun 20:00 → local Mon 05:00
    const result = utcDayHourToLocal(7, 20, "Asia/Tokyo");
    expect(result.localHour).toBe(5);
    expect(result.localDay).toBe(1); // Monday
  });
});
