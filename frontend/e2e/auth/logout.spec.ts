import { test, expect } from "../fixtures";

test.describe("Logout", () => {
  test("user can log out and cannot access dashboard", async ({
    authedPage,
  }) => {
    await authedPage.goto("/dashboard");

    await authedPage.getByTestId("nav-user-trigger").click();
    await authedPage.getByTestId("logout-menu-item").click();

    await expect(
      authedPage.getByText("Logged out", { exact: true }),
    ).toBeVisible({ timeout: 5_000 });
    await expect(authedPage).not.toHaveURL(/\/dashboard/);

    await authedPage.context().clearCookies();
    await authedPage.goto("/dashboard");
    await authedPage.waitForURL(/\/login/, { timeout: 10_000 });
    await expect(authedPage).toHaveURL(/\/login/);
  });
});
