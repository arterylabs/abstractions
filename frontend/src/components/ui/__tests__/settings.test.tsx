import { useState } from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { DeleteRepositoryDialog } from "@/app/(dashboard)/repositories/[id]/settings/DeleteRepositoryDialog";
import { utcDayHourToLocal } from "@/lib/timezone";

describe("DeleteRepositoryDialog", () => {
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
  // DST-free zones so expectations are stable year-round.
  it.each([
    {
      name: "positive offset without crossing midnight",
      // Asia/Tokyo = UTC+9: UTC Wed 10:00 → local Wed 19:00
      day: 3,
      hour: 10,
      tz: "Asia/Tokyo",
      localHour: 19,
      localDay: 3,
    },
    {
      name: "negative offset without crossing midnight",
      // Pacific/Honolulu = UTC-10: UTC Wed 15:00 → local Wed 05:00
      day: 3,
      hour: 15,
      tz: "Pacific/Honolulu",
      localHour: 5,
      localDay: 3,
    },
    {
      name: "rolls forward across midnight",
      // Asia/Tokyo: UTC Wed 20:00 → local Thu 05:00
      day: 3,
      hour: 20,
      tz: "Asia/Tokyo",
      localHour: 5,
      localDay: 4,
    },
    {
      name: "rolls back across midnight",
      // Pacific/Honolulu: UTC Mon 05:00 → local Sun 19:00
      day: 1,
      hour: 5,
      tz: "Pacific/Honolulu",
      localHour: 19,
      localDay: 7,
    },
    {
      name: "Sunday→Monday wrap",
      // Asia/Tokyo: UTC Sun 20:00 → local Mon 05:00
      day: 7,
      hour: 20,
      tz: "Asia/Tokyo",
      localHour: 5,
      localDay: 1,
    },
  ])("$name", ({ day, hour, tz, localHour, localDay }) => {
    const result = utcDayHourToLocal(day, hour, tz);
    expect(result.localHour).toBe(localHour);
    expect(result.localDay).toBe(localDay);
  });
});
