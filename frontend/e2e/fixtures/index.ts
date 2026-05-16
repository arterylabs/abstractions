import { test as base, Page } from "@playwright/test";
import { signup } from "../utils/auth";

type TestFixtures = {
  authedPage: Page;
  user: {
    name: string;
    email: string;
    password: string;
  };
};

export const test = base.extend<TestFixtures>({
  user: async ({}, setFixture, testInfo) => {
    const suffix = `${Date.now()}-${testInfo.parallelIndex}-${Math.floor(Math.random() * 10000)}`;
    await setFixture({
      name: "E2E User",
      email: `e2e-${suffix}@test.arterylabs.com`,
      password: "Password123!",
    });
  },

  authedPage: async ({ page, user }, setFixture) => {
    await signup(page, user.name, user.email, user.password);

    // New users have no workspace, so WorkspaceContext redirects to /workspaces/create.
    // Wait briefly to see if a redirect happens, then complete workspace creation.
    try {
      await page.waitForURL(/\/workspaces\/create/, { timeout: 5_000 });
      const workspaceName = `${user.name}'s Workspace`;
      await page.getByTestId("create-workspace-name-input").fill(workspaceName);
      await page.getByTestId("create-workspace-submit").click();
      await page.waitForURL(/\/dashboard/, { timeout: 15_000 });
    } catch {
      // Already on /dashboard with an existing workspace — nothing to do.
    }

    await setFixture(page);
  },
});

export { expect } from "@playwright/test";
