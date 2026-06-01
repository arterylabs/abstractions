import { test, expect } from "@playwright/test";
import { signup } from "../utils/auth";

const apiURL =
  process.env.PLAYWRIGHT_API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:1234";

async function verifyTestUser(email: string): Promise<void> {
  await fetch(`${apiURL}/test/verify_user`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

test.describe("Email verification", () => {
  test("shows email verification modal after signup", async ({ page }) => {
    const email = `verify-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Verify Test User", email, "Password123!");

    const modal = page.getByTestId("email-verification-modal");
    await expect(modal).toBeVisible();
    await expect(modal).toContainText(email);
    await expect(modal).toContainText("Verify your email address");
  });

  test("resend button sends a verification email and shows cooldown", async ({
    page,
  }) => {
    const email = `verify-resend-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Verify Resend User", email, "Password123!");

    const resendButton = page.getByTestId("resend-verification-button");
    await expect(resendButton).toBeVisible();
    await expect(resendButton).toBeEnabled();

    await resendButton.click();

    // After sending, the button should show a cooldown timer
    await expect(resendButton).toHaveText(/Resend in \d+s/, { timeout: 5_000 });
    await expect(resendButton).toBeDisabled();
  });

  test("modal disappears after account is verified", async ({ page }) => {
    const email = `verify-disappear-${Date.now()}@test.arterylabs.com`;
    await signup(page, "Verify Disappear User", email, "Password123!");

    await expect(page.getByTestId("email-verification-modal")).toBeVisible();

    await verifyTestUser(email);

    await page.reload();
    await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });

    await expect(
      page.getByTestId("email-verification-modal"),
    ).not.toBeVisible();
  });
});
