import { test, expect } from "../fixtures";

test.describe("Dashboard navigation", () => {
  test("navigation between dashboard sections works", async ({
    authedPage,
  }) => {
    await authedPage.goto("/dashboard");

    await expect(authedPage.getByTestId("dashboard-nav")).toBeVisible();
    await expect(authedPage.getByTestId("nav-dashboard-link")).toBeVisible();
    await expect(authedPage.getByTestId("nav-settings-link")).toBeVisible();

    await authedPage.getByTestId("nav-settings-link").click();
    await expect(authedPage).toHaveURL(/\/settings/);

    await authedPage.getByTestId("nav-dashboard-link").click();
    await expect(authedPage).toHaveURL(/\/repositories/);
  });

  test("sidebar navigates to getting started, integrations, and AI", async ({
    authedPage,
  }) => {
    await authedPage.goto("/repositories");
    await expect(authedPage.getByTestId("dashboard-nav")).toBeVisible();

    await authedPage.getByTestId("nav-getting_started-link").click();
    await expect(authedPage).toHaveURL(/\/getting_started/);
    await expect(
      authedPage.getByRole("heading", { name: "Getting started" }),
    ).toBeVisible();

    await authedPage.getByTestId("nav-integrations-link").click();
    await expect(authedPage).toHaveURL(/\/integrations/);
    await expect(
      authedPage.getByRole("heading", { name: "Integrations" }),
    ).toBeVisible();

    await authedPage.getByTestId("nav-ai-link").click();
    await expect(authedPage).toHaveURL(/\/ai/);
    await expect(
      authedPage.getByRole("heading", { name: "AI", exact: true, level: 1 }),
    ).toBeVisible();
  });
});
