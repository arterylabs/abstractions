import { test, expect } from "../fixtures";

test.describe("Dashboard navigation", () => {
  test("navigation between dashboard sections works", async ({
    authedPage,
  }) => {
    await authedPage.goto("/dashboard");

    await expect(authedPage.getByTestId("dashboard-nav")).toBeVisible();
    await expect(authedPage.getByTestId("nav-dashboard-link")).toBeVisible();
    await expect(authedPage.getByTestId("nav-settings-link")).toBeVisible();

    await authedPage.getByTestId("nav-settings-link").click();
    await expect(authedPage).toHaveURL(/\/settings/);

    await authedPage.getByTestId("nav-dashboard-link").click();
    await expect(authedPage).toHaveURL(/\/repositories/);
  });
});
