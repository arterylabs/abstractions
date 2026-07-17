import { test, expect } from "@playwright/test";
import { fillField, login, signup, uniqueTestEmail } from "../utils/auth";

test.describe("Mobile auth", () => {
  test("user can sign up on mobile", async ({ page }) => {
    const email = uniqueTestEmail("mobile-signup");
    await signup(page, "Mobile Test User", email, "Password123!");

    await expect(page).toHaveURL(/\/(repositories|dashboard|verify)/);
  });

  test("login form is usable on mobile", async ({ page }) => {
    const email = uniqueTestEmail("mobile-login");
    await signup(page, "Mobile Login User", email, "Password123!");
    await page.context().clearCookies();

    await login(page, email, "Password123!");
    // /dashboard server-redirects to /repositories
    await expect(page).toHaveURL(/\/(repositories|dashboard)/);
  });

  test("login shows error for invalid credentials on mobile", async ({
    page,
  }) => {
    await page.goto("/login");

    // Use fillField — raw fill() is flaky on WebKit for controlled inputs.
    await fillField(page, "login-email-input", "nope@test.arterylabs.com");
    await fillField(page, "login-password-input", "wrongpassword");
    await page.getByTestId("login-submit-button").click();

    await expect(
      page.getByText("Error logging in", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });
  });
});
