import { test, expect } from "@playwright/test";
import { login, signup, uniqueTestEmail } from "../utils/auth";

test.describe("Login", () => {
  test("user can log in with valid credentials", async ({ page }) => {
    const email = uniqueTestEmail("login");
    const password = "Password123!";

    await signup(page, "Playwright Login User", email, password);
    await page.context().clearCookies();
    await login(page, email, password);

    await expect(page).toHaveURL(/\/(repositories|dashboard)/);
  });

  test("shows error for invalid credentials", async ({ page }) => {
    await page.goto("/login");

    await page
      .getByTestId("login-email-input")
      .fill("does-not-exist@test.arterylabs.com");
    await page.getByTestId("login-password-input").fill("wrong-password");
    await page.getByTestId("login-submit-button").click();

    await expect(
      page.getByText("Error logging in", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });
  });
});
