import { test, expect } from "@playwright/test";

test.describe("Forgot password", () => {
  test("user can request a password reset email", async ({ page }) => {
    await page.goto("/forgot");

    await page
      .getByTestId("forgot-email-input")
      .fill(`reset-${Date.now()}@test.arterylabs.com`);
    await page.getByTestId("forgot-submit-button").click();

    // After a successful request the page transitions to the "check inbox" state
    await expect(page.getByText("Check Your Inbox")).toBeVisible({
      timeout: 10_000,
    });
  });

  test("forgot password page renders the email form", async ({ page }) => {
    await page.goto("/forgot");

    await expect(page.getByTestId("forgot-email-input")).toBeVisible();
    await expect(page.getByTestId("forgot-submit-button")).toBeVisible();
    await expect(
      page.getByText("Forgot Password?", { exact: false }),
    ).toBeVisible();
  });

  test("forgot password page has a link back to login", async ({ page }) => {
    await page.goto("/forgot");

    await page.getByRole("link", { name: /login/i }).click();

    await expect(page).toHaveURL(/\/login/);
  });
});
