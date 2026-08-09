import { test, expect } from "../fixtures";

test.describe("Repositories", () => {
  test("repositories page shows empty state for new workspace", async ({
    authedPage,
  }) => {
    await authedPage.goto("/repositories");

    await expect(
      authedPage.getByRole("heading", { name: "Repositories" }),
    ).toBeVisible();
    await expect(
      authedPage.getByText("Connect your first repository"),
    ).toBeVisible({ timeout: 15_000 });
    await expect(
      authedPage.getByRole("button", { name: /Connect repository/i }).first(),
    ).toBeVisible();
  });

  test("connect repository button opens connect modal", async ({
    authedPage,
  }) => {
    await authedPage.goto("/repositories");

    // Wait for the empty state so the page has finished loading.
    await expect(
      authedPage.getByText("Connect your first repository"),
    ).toBeVisible({ timeout: 15_000 });

    await authedPage
      .getByRole("button", { name: /Connect repository/i })
      .first()
      .click();

    // Modal opens in a "checking" state, then settles on the OAuth prompt for
    // workspaces without a GitHub connection.
    const dialog = authedPage.getByRole("dialog");
    await expect(dialog).toBeVisible();
    await expect(
      dialog.getByRole("heading", {
        name: /Connect a (Repository|GitHub Repository)/i,
      }),
    ).toBeVisible();

    // Fresh workspaces should land on the OAuth step.
    await expect(
      dialog.getByRole("button", { name: /Connect to GitHub/i }),
    ).toBeVisible({ timeout: 15_000 });
  });
});
