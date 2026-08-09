import { test, expect } from "../fixtures";
import { createWorkspace } from "../utils/auth";

test.describe("Workspace management", () => {
  test("user can create a workspace and switch between workspaces", async ({
    authedPage,
  }) => {
    const secondWorkspaceName = `E2E Workspace ${Date.now()}`;

    const initialWorkspaceName = await authedPage
      .getByTestId("workspace-switcher-name")
      .textContent();

    await createWorkspace(authedPage, secondWorkspaceName);

    await authedPage.reload();
    await authedPage.waitForURL(/\/repositories/, { timeout: 15_000 });

    await expect(
      authedPage.getByTestId("workspace-switcher-trigger"),
    ).toContainText(secondWorkspaceName);

    await authedPage.getByTestId("workspace-switcher-trigger").click();
    await expect(authedPage.getByText(secondWorkspaceName)).toBeVisible();

    const firstWorkspaceItem = authedPage.getByRole("menuitem").filter({
      hasText: initialWorkspaceName?.trim() ?? "",
    });
    await firstWorkspaceItem.click();

    await expect(
      authedPage.getByTestId("workspace-switcher-trigger"),
    ).not.toContainText(secondWorkspaceName);
  });

  test("user can update workspace name in workspace settings", async ({
    authedPage,
  }) => {
    const updatedName = `Updated WS ${Date.now()}`;

    await authedPage.goto("/settings/workspace");

    const nameInput = authedPage.getByTestId("workspace-name-input");
    await nameInput.clear();
    await nameInput.fill(updatedName);

    await authedPage.getByTestId("workspace-save-button").click();

    await expect(
      authedPage.getByText("Settings saved", { exact: false }),
    ).toBeVisible();

    await authedPage.reload();
    await expect(authedPage.getByTestId("workspace-name-input")).toHaveValue(
      updatedName,
    );
  });
});
