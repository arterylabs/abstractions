import { test, expect } from "@playwright/test";

test.describe("Reset password", () => {
  test("reset password page renders without token", async ({ page }) => {
    await page.goto("/reset");

    // CardTitle is a div (not a heading role) in this design system.
    await expect(page.getByText("Reset Password").first()).toBeVisible();
    await expect(page.getByText("Set your new password here")).toBeVisible();
    await expect(page.getByLabel("Password", { exact: true })).toBeVisible();
    await expect(page.getByLabel("Confirm Password")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Reset Password" }),
    ).toBeVisible();
    await expect(page.getByRole("link", { name: /login/i })).toBeVisible();
  });
});
