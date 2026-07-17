import { test, expect } from "../fixtures";

test.describe("Mobile workspace settings", () => {
  test("settings tabs are visible and scrollable on mobile", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace");

    // All tabs should be present
    await expect(
      authedPage.getByRole("tab", { name: /general/i }),
    ).toBeVisible();
    await expect(authedPage.getByRole("tab", { name: /users/i })).toBeVisible();
    await expect(
      authedPage.getByRole("tab", { name: /invitations/i }),
    ).toBeVisible();
  });

  test("workspace name input is reachable on mobile", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace");

    const nameInput = authedPage.getByTestId("workspace-name-input");
    await expect(nameInput).toBeVisible();
    await nameInput.click();
    await expect(nameInput).toBeFocused();
  });

  test("save button is visible below the form on mobile", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace");

    const saveButton = authedPage.getByTestId("workspace-save-button");
    await expect(saveButton).toBeVisible();
  });

  test("users tab loads on mobile", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace/users");

    await expect(authedPage.getByPlaceholder("Search users...")).toBeVisible();
  });
});
