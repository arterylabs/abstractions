import { test, expect } from "../fixtures";

test.describe("Mobile smoke", () => {
  test("can open sidebar and navigate on mobile", async ({ authedPage }) => {
    await authedPage.goto("/repositories");

    const sidebarTrigger = authedPage.getByRole("button", {
      name: /toggle sidebar/i,
    });
    await expect(sidebarTrigger).toBeVisible();

    await sidebarTrigger.click();
    await expect(authedPage.getByTestId("dashboard-nav")).toBeVisible();

    await authedPage.getByTestId("nav-settings-link").click();
    await expect(authedPage).toHaveURL(/\/settings/);
  });
});
