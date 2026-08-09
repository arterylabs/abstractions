import { test, expect } from "../fixtures";

test.describe("Advanced account settings", () => {
  test("delete account is disabled until DELETE is typed", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/advanced");

    await expect(
      authedPage.getByRole("heading", { name: "Delete Account" }),
    ).toBeVisible();

    await authedPage.getByRole("button", { name: /delete account/i }).click();

    await expect(
      authedPage.getByText("Delete your account permanently?"),
    ).toBeVisible();

    const confirmButton = authedPage.getByRole("button", {
      name: /delete account permanently/i,
    });
    await expect(confirmButton).toBeDisabled();

    await authedPage.getByPlaceholder("Type DELETE to confirm").fill("DELETE");
    await expect(confirmButton).toBeEnabled();
  });
});
