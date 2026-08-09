import { test, expect } from "../fixtures";

/**
 * Full connect journey: mock GitHub repo listing (no real OAuth token),
 * then POST a real repository row via the API and assert it appears in the UI.
 */
test.describe("Repositories", () => {
  test("user can connect a repository from the connect flow", async ({
    authedPage,
  }) => {
    const fakeRepo = {
      id: 424242,
      fullName: "acme/e2e-connected-repo",
      defaultBranch: "main",
      private: false,
      description: "E2E connected repository",
      htmlUrl: "https://github.com/acme/e2e-connected-repo",
    };

    // Mock GitHub list so the connect page shows repos without OAuth.
    await authedPage.route("**/github/repos", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([fakeRepo]),
      });
    });

    await authedPage.goto("/repositories/connect");

    await expect(
      authedPage.getByText(fakeRepo.fullName, { exact: false }),
    ).toBeVisible({ timeout: 15_000 });

    await authedPage.getByRole("button", { name: "Connect" }).first().click();

    // Success toast + navigation to the post-connect route prove the real
    // POST /repositories call succeeded (list uses mocked GitHub only).
    await expect(
      authedPage.getByText(/Repository connected|connected successfully/i),
    ).toBeVisible({ timeout: 15_000 });
    await expect(authedPage).toHaveURL(/\/repositories\/connect\/[a-f0-9-]+/, {
      timeout: 15_000,
    });

    await authedPage.goto("/repositories");
    // List may render `name` or `fullName` depending on API field mapping.
    await expect(
      authedPage.getByText(fakeRepo.fullName, { exact: false }),
    ).toBeVisible({ timeout: 15_000 });
  });
});
