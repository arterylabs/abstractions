import { test, expect } from "@playwright/test";
import { signup, verifyTestUser } from "../utils/auth";

test.describe("Email verification", () => {
  test("modal appears after signup, supports resend, and clears when verified", async ({
    page,
  }) => {
    const email = `verify-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Verify Test User", email, "Password123!");

    const modal = page.getByTestId("email-verification-modal");
    await expect(modal).toBeVisible();
    await expect(modal).toContainText(email);

    const resendButton = page.getByTestId("resend-verification-button");
    await expect(resendButton).toBeEnabled();
    await resendButton.click();
    await expect(resendButton).toHaveText(/Resend in \d+s/, { timeout: 5_000 });
    await expect(resendButton).toBeDisabled();

    await verifyTestUser(email);
    await page.reload();
    await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });
    await expect(
      page.getByTestId("email-verification-modal"),
    ).not.toBeVisible();
  });
});
