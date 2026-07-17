import { test, expect } from "../fixtures";

test.describe("Advanced account settings", () => {
  test("advanced settings page renders the delete account section", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/advanced");

    // The heading should be present on the page
    await expect(
      authedPage.getByRole("heading", { name: "Delete Account" }),
    ).toBeVisible();
  });

  test("delete account button opens confirmation dialog", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/advanced");

    await authedPage.getByRole("button", { name: /delete account/i }).click();

    await expect(
      authedPage.getByText("Delete your account permanently?"),
    ).toBeVisible();
  });

  test("delete account button is disabled until DELETE is typed", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/advanced");

    await authedPage.getByRole("button", { name: /delete account/i }).click();

    const confirmButton = authedPage.getByRole("button", {
      name: /delete account permanently/i,
    });

    // Button should be disabled initially
    await expect(confirmButton).toBeDisabled();

    // Type the confirmation text
    await authedPage.getByPlaceholder("Type DELETE to confirm").fill("DELETE");

    // Button should now be enabled
    await expect(confirmButton).toBeEnabled();
  });
});
