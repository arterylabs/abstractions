import { test as base, Page } from "@playwright/test";
import {
  signup,
  verifyTestUser,
  waitForAuthenticatedShell,
} from "../utils/auth";

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

    // Mark the freshly-signed-up test account as verified so the
    // email-verification modal does not block test interactions and
    // the backend allows verification-gated actions (workspace creation, etc.).
    await verifyTestUser(user.email);

    // Land on /repositories explicitly so the /dashboard → /repositories
    // server redirect cannot race later page.goto() calls (especially WebKit).
    // Avoid racing waitForResponse: /users/me may complete before the waiter
    // is registered under load. waitForAuthenticatedShell is the source of truth.
    await page.goto("/repositories", { waitUntil: "domcontentloaded" });
    await waitForAuthenticatedShell(page);

    await setFixture(page);
  },
});

export { expect } from "@playwright/test";
