import { test, expect } from "@playwright/test";
import { login, signup } from "../utils/auth";

test.describe("Mobile auth", () => {
  test("user can sign up on mobile", async ({ page }) => {
    const email = `mobile-signup-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Mobile Test User", email, "Password123!");

    await expect(page).toHaveURL(/\/(repositories|dashboard|verify)/);
  });

  test("login form is usable on mobile", async ({ page }) => {
    const email = `mobile-login-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Mobile Login User", email, "Password123!");
    await page.context().clearCookies();

    await login(page, email, "Password123!");
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test("login shows error for invalid credentials on mobile", async ({
    page,
  }) => {
    await page.goto("/login");

    await page
      .getByTestId("login-email-input")
      .fill("nope@test.arterylabs.com");
    await page.getByTestId("login-password-input").fill("wrongpassword");
    await page.getByTestId("login-submit-button").click();

    await expect(
      page.getByText("Error logging in", { exact: false }),
    ).toBeVisible();
  });
});
