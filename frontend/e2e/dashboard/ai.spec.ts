import { test, expect } from "../fixtures";

test.describe("AI settings", () => {
  test("AI settings page renders providers section", async ({ authedPage }) => {
    await authedPage.goto("/ai");

    await expect(
      authedPage.getByRole("heading", { name: "AI", exact: true, level: 1 }),
    ).toBeVisible();

    await expect(
      authedPage.getByText(
        "Configure the models and API keys Abstractions uses to index your code and write your newsletters.",
      ),
    ).toBeVisible();

    // Section marks use "Providers" / "Models" labels
    await expect(authedPage.getByText("Providers").first()).toBeVisible({
      timeout: 15_000,
    });
    await expect(authedPage.getByText("Models").first()).toBeVisible();

    // At least one provider card should load from the platform catalog
    await expect(
      authedPage.getByText(/Anthropic|OpenAI|xAI|Grok/i).first(),
    ).toBeVisible({ timeout: 15_000 });
  });
});
