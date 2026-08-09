import { test, expect } from "@playwright/test";

test.describe("Protected routes", () => {
  test("unauthenticated user is redirected away from dashboard routes", async ({
    page,
  }) => {
    await page.goto("/repositories");
    await page.waitForURL(/\/login/, { timeout: 15_000 });
    await expect(page).toHaveURL(/\/login/);

    await page.goto("/settings/workspace");
    await page.waitForURL(/\/login/, { timeout: 15_000 });
    await expect(page).toHaveURL(/\/login/);
  });
});
