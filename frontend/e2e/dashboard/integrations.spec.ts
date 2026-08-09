import { test, expect } from "../fixtures";

test.describe("Integrations", () => {
  test("integrations page renders for workspace", async ({ authedPage }) => {
    await authedPage.goto("/integrations");

    await expect(
      authedPage.getByRole("heading", { name: "Integrations" }),
    ).toBeVisible();
    await expect(
      authedPage.getByText(
        "Connect the tools your team already uses so Abstractions can pull in the right context.",
      ),
    ).toBeVisible();

    // Supported types load from the API — GitHub is always registered.
    // CardTitle is a div, so match text rather than a heading role.
    await expect(authedPage.getByText("GitHub", { exact: true })).toBeVisible({
      timeout: 15_000,
    });
    await expect(
      authedPage.getByText(
        "Connect your GitHub account to let Abstractions access your repositories.",
      ),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("button", { name: /Connect/i }).first(),
    ).toBeVisible();
  });
});
