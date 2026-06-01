import { test, expect } from "../fixtures";

test.describe("Mobile navigation", () => {
  test("sidebar is hidden and togglable via trigger on mobile", async ({
    authedPage,
  }) => {
    await authedPage.goto("/repositories");

    // On mobile the sidebar is off-canvas; the toggle trigger should be visible
    const sidebarTrigger = authedPage.getByRole("button", {
      name: /toggle sidebar/i,
    });
    await expect(sidebarTrigger).toBeVisible();

    // The sidebar nav should not be visible before opening
    const nav = authedPage.getByTestId("dashboard-nav");
    await expect(nav).not.toBeVisible();

    // Open the sidebar
    await sidebarTrigger.click();
    await expect(nav).toBeVisible();
  });

  test("can navigate to settings via mobile sidebar", async ({
    authedPage,
  }) => {
    await authedPage.goto("/repositories");

    const sidebarTrigger = authedPage.getByRole("button", {
      name: /toggle sidebar/i,
    });
    await sidebarTrigger.click();

    await authedPage.getByTestId("nav-settings-link").click();
    await expect(authedPage).toHaveURL(/\/settings/);
  });

  test("can navigate to repositories via mobile sidebar", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace");

    const sidebarTrigger = authedPage.getByRole("button", {
      name: /toggle sidebar/i,
    });
    await sidebarTrigger.click();

    await authedPage.getByTestId("nav-dashboard-link").click();
    await expect(authedPage).toHaveURL(/\/repositories/);
  });
});
