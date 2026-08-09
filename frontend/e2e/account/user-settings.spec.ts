import { test, expect } from "../fixtures";

test.describe("User account settings", () => {
  test("user can update their display name", async ({ authedPage }) => {
    const updatedName = `Updated Name ${Date.now()}`;

    await authedPage.goto("/account");

    const nameInput = authedPage.getByTestId("account-name-input");
    await nameInput.clear();
    await nameInput.fill(updatedName);

    await authedPage.getByTestId("account-save-button").click();

    await expect(
      authedPage.getByText("Profile saved", { exact: false }),
    ).toBeVisible();

    await authedPage.reload();
    await expect(authedPage.getByTestId("account-name-input")).toHaveValue(
      updatedName,
    );
  });
});
