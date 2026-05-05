import { test, expect } from "@playwright/test";
import { login, signup } from "../utils/auth";

test.describe("Login", () => {
  test("user can log in with valid credentials", async ({ page }) => {
    const email = `login-${Date.now()}@example.com`;
    const password = "Password123!";

    await signup(page, "Playwright Login User", email, password);
    await contextSignOut(page);
    await login(page, email, password);

    await expect(page).toHaveURL(/\/dashboard/);
  });

  test("shows error for invalid credentials", async ({ page }) => {
    await page.goto("/login");

    await page
      .getByTestId("login-email-input")
      .fill("does-not-exist@example.com");
    await page.getByTestId("login-password-input").fill("wrong-password");
    await page.getByTestId("login-submit-button").click();

    await expect(
      page.getByText("Error logging in", { exact: false }),
    ).toBeVisible();
  });
});

async function contextSignOut(page: import("@playwright/test").Page) {
  await page.context().clearCookies();
}
