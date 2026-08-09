import { test, expect } from "@playwright/test";

test.describe("Forgot password", () => {
  test("user can request a password reset email", async ({ page }) => {
    await page.goto("/forgot");

    await page
      .getByTestId("forgot-email-input")
      .fill(`reset-${Date.now()}@test.arterylabs.com`);
    await page.getByTestId("forgot-submit-button").click();

    await expect(page.getByText("Check Your Inbox")).toBeVisible({
      timeout: 10_000,
    });
  });
});
