import { test, expect } from "@playwright/test";
import { login, signup, uniqueTestEmail } from "../utils/auth";

test.describe("Login", () => {
  test("user can log in with valid credentials", async ({ page }) => {
    const email = uniqueTestEmail("login");
    const password = "Password123!";

    await signup(page, "Playwright Login User", email, password);
    await page.context().clearCookies();
    await login(page, email, password);

    // /dashboard server-redirects to /repositories
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

  test("login page links to signup and forgot password", async ({ page }) => {
    await page.goto("/login");

    await expect(page.getByRole("link", { name: /sign up/i })).toHaveAttribute(
      "href",
      "/signup",
    );
    await expect(
      page.getByRole("link", { name: /forgot password/i }),
    ).toHaveAttribute("href", "/forgot");

    await page.getByRole("link", { name: /sign up/i }).click();
    await expect(page).toHaveURL(/\/signup/);

    await page.goto("/login");
    await page.getByRole("link", { name: /forgot password/i }).click();
    await expect(page).toHaveURL(/\/forgot/);
  });
});
