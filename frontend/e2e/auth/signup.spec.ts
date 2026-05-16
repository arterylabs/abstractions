import { test, expect } from "@playwright/test";
import { signup } from "../utils/auth";

test.describe("Signup", () => {
  test("user can sign up with valid details", async ({ page }) => {
    const email = `signup-${Date.now()}@test.arterylabs.com`;

    await signup(page, "Playwright Signup User", email, "Password123!");

    await expect(page).toHaveURL(/\/dashboard/);
  });

  test("shows an error for duplicate email signup", async ({ page }) => {
    const email = `signup-duplicate-${Date.now()}@test.arterylabs.com`;
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
    ).toBeVisible();
  });
});
