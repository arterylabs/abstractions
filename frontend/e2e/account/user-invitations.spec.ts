import { test, expect } from "../fixtures";

test.describe("User invitations inbox", () => {
  test("pending invitations inbox shows empty state", async ({
    authedPage,
  }) => {
    await authedPage.goto("/me/invitations");

    await expect(authedPage.getByText("No Pending Invitations")).toBeVisible({
      timeout: 15_000,
    });
    await expect(
      authedPage.getByText(
        "You don't have any pending invitations at the moment.",
      ),
    ).toBeVisible();
  });
});
