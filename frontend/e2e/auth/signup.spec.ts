import { test, expect } from "@playwright/test";
import { signup, uniqueTestEmail } from "../utils/auth";

test.describe("Signup", () => {
  test("user can sign up with valid details", async ({ page }) => {
    const email = uniqueTestEmail("signup");

    await signup(page, "Playwright Signup User", email, "Password123!");

    // /dashboard server-redirects to /repositories
    await expect(page).toHaveURL(/\/(repositories|dashboard)/);
  });

  test("shows an error for duplicate email signup", async ({ page }) => {
    const email = uniqueTestEmail("signup-duplicate");
    const password = "Password123!";

    await signup(page, "Duplicate Signup User", email, password);
    await page.context().clearCookies();

    await page.goto("/signup");
    await page.getByTestId("signup-name-input").fill("Duplicate Signup User");
    await page.getByTestId("signup-email-input").fill(email);
    await page.getByTestId("signup-password-input").fill(password);
    await page.getByTestId("signup-submit-button").click();

    await expect(
      page.getByText("Error signing up", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });
  });
});
