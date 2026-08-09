import { test, expect } from "../fixtures";

test.describe("Getting started", () => {
  test("getting started page lists onboarding steps", async ({
    authedPage,
  }) => {
    await authedPage.getByTestId("nav-getting_started-link").click();

    await expect(authedPage).toHaveURL(/\/getting_started/);
    await expect(
      authedPage.getByRole("heading", { name: "Getting started" }),
    ).toBeVisible();

    // Three setup steps for a fresh workspace
    await expect(
      authedPage.getByRole("heading", { name: "Connect your AI providers" }),
    ).toBeVisible({ timeout: 15_000 });
    await expect(
      authedPage.getByRole("heading", {
        name: "Connect your GitHub repository",
      }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("heading", { name: "Send your first newsletter" }),
    ).toBeVisible();

    await expect(authedPage.getByText(/of 3 complete/i)).toBeVisible();
  });
});
