import { test, expect } from "../fixtures";

test.describe("Workspace management", () => {
  test("user can create a new workspace", async ({ authedPage }) => {
    const workspaceName = `E2E Workspace ${Date.now()}`;

    await authedPage.getByTestId("account-switcher-trigger").click();
    await authedPage.getByTestId("create-workspace-link").click();

    await expect(authedPage).toHaveURL(/\/workspaces\/create/);
    await expect(
      authedPage.getByRole("heading", { name: "Create your workspace" }),
    ).toBeVisible();

    await authedPage
      .getByTestId("create-workspace-name-input")
      .fill(workspaceName);
    await authedPage.getByTestId("create-workspace-submit").click();

    await authedPage.waitForURL(/\/dashboard/, { timeout: 15_000 });
    await expect(authedPage.getByText(workspaceName)).toBeVisible();
  });

  test("newly created workspace appears in the workspace switcher", async ({
    authedPage,
  }) => {
    const workspaceName = `Switcher Workspace ${Date.now()}`;

    await authedPage.getByTestId("account-switcher-trigger").click();
    await authedPage.getByTestId("create-workspace-link").click();
    await authedPage
      .getByTestId("create-workspace-name-input")
      .fill(workspaceName);
    await authedPage.getByTestId("create-workspace-submit").click();
    await authedPage.waitForURL(/\/dashboard/, { timeout: 15_000 });

    // Open the switcher and confirm the new workspace is listed
    await authedPage.getByTestId("account-switcher-trigger").click();
    await expect(authedPage.getByText(workspaceName)).toBeVisible();
  });

  test("user can switch between workspaces", async ({ authedPage }) => {
    const secondWorkspaceName = `Second WS ${Date.now()}`;

    // Capture the name of the initial workspace shown in the header
    await authedPage.getByTestId("account-switcher-trigger").click();
    const initialWorkspaceName = await authedPage
      .getByTestId("account-switcher-trigger")
      .textContent();

    // Create a second workspace
    await authedPage.getByTestId("create-workspace-link").click();
    await authedPage
      .getByTestId("create-workspace-name-input")
      .fill(secondWorkspaceName);
    await authedPage.getByTestId("create-workspace-submit").click();
    await authedPage.waitForURL(/\/dashboard/, { timeout: 15_000 });

    // Confirm the switcher now shows the new workspace
    await expect(
      authedPage.getByTestId("account-switcher-trigger"),
    ).toContainText(secondWorkspaceName);

    // Switch back to the first workspace via the dropdown
    await authedPage.getByTestId("account-switcher-trigger").click();
    const firstWorkspaceItem = authedPage.getByRole("menuitem").filter({
      hasText: initialWorkspaceName?.trim() ?? "",
    });
    await firstWorkspaceItem.click();

    // Switcher header should reflect the first workspace again
    await expect(
      authedPage.getByTestId("account-switcher-trigger"),
    ).not.toContainText(secondWorkspaceName);
  });

  test("user can navigate to workspace settings from the switcher", async ({
    authedPage,
  }) => {
    await authedPage.getByTestId("account-switcher-trigger").click();
    await authedPage.getByTestId("workspace-settings-link").click();

    await expect(authedPage).toHaveURL(/\/settings\/workspace/);
    await expect(
      authedPage.getByRole("heading", { name: "Workspace Settings" }),
    ).toBeVisible();
  });

  test("user can update workspace name in workspace settings", async ({
    authedPage,
    user,
  }) => {
    const updatedName = `Updated WS ${Date.now()}`;

    await authedPage.goto("/settings/workspace");

    const nameInput = authedPage.getByTestId("workspace-name-input");
    await nameInput.clear();
    await nameInput.fill(updatedName);

    await authedPage.getByTestId("workspace-save-button").click();

    await expect(
      authedPage.getByText("Settings updated", { exact: false }),
    ).toBeVisible();

    // Reload and confirm the name persisted
    await authedPage.reload();
    await expect(authedPage.getByTestId("workspace-name-input")).toHaveValue(
      updatedName,
    );

    // User fixture keeps name in scope — suppress TS "unused" warnings
    void user;
  });

  test("workspace settings shows billing email field", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace");

    await expect(
      authedPage.getByTestId("workspace-billing-email-input"),
    ).toBeVisible();
    await expect(authedPage.getByTestId("workspace-save-button")).toBeVisible();
  });

  test("workspace settings shows all settings tabs", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace");

    await expect(
      authedPage.getByRole("link", { name: /General/i }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("link", { name: /Users/i }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("link", { name: /Invitations/i }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("link", { name: /Billing/i }),
    ).toBeVisible();
  });
});
