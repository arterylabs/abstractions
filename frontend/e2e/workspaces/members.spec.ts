import { test, expect } from "../fixtures";

test.describe("Workspace members", () => {
  test("workspace members page is accessible from workspace settings", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace/users");

    // Page should load without errors
    await expect(authedPage).toHaveURL(/\/settings\/workspace\/users/);
  });

  test("workspace members page lists the current user", async ({
    authedPage,
    user,
  }) => {
    await authedPage.goto("/settings/workspace/users");

    // The authenticated user should appear in the member list
    await expect(
      authedPage.getByLabel("Infinite scroll content").getByText(user.name),
    ).toBeVisible({
      timeout: 10_000,
    });
  });

  test("workspace members page shows search bar", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace/users");

    // Search bar should be present on the page
    await expect(authedPage.getByPlaceholder("Search users...")).toBeVisible();
  });
});
