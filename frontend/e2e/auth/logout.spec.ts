import { test, expect } from "../fixtures";

test.describe("Logout", () => {
  test("user can log out", async ({ authedPage }) => {
    await authedPage.goto("/dashboard");

    // Open the user menu in the sidebar nav
    await authedPage.getByTestId("nav-user-trigger").click();

    await authedPage.getByTestId("logout-menu-item").click();

    // A success toast should confirm the logout
    await expect(
      authedPage.getByText("Logged out", { exact: true }),
    ).toBeVisible({ timeout: 5_000 });

    // Logout redirects to the landing page (/) — confirm we left the dashboard
    await expect(authedPage).not.toHaveURL(/\/dashboard/);
  });

  test("logged out user cannot access dashboard", async ({ authedPage }) => {
    // Clear cookies to simulate being logged out
    await authedPage.context().clearCookies();

    await authedPage.goto("/dashboard");

    // Should be redirected away from dashboard (to login or landing)
    await authedPage.waitForURL(/\/login/, { timeout: 10_000 });
    await expect(authedPage).toHaveURL(/\/login/);
  });
});
