import { test, expect } from "@playwright/test";

test.describe("Documentation", () => {
  test("docs index page loads", async ({ page }) => {
    await page.goto("/docs");

    await expect(
      page.getByRole("heading", { name: "Introduction" }),
    ).toBeVisible({ timeout: 15_000 });
    await expect(
      page.getByText(
        "Abstractions turns your GitHub activity into automated developer newsletters",
        { exact: false },
      ),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "What you can do" }),
    ).toBeVisible();
  });
});
