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

const apiURL =
  process.env.PLAYWRIGHT_API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:1234";

/** Mark a single test account as verified via the dev-only endpoint. */
async function verifyTestUser(email: string): Promise<void> {
  await fetch(`${apiURL}/test/verify-user`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

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

    // Reload so the frontend re-fetches /users/me and sees verifiedAt.
    await page.reload();
    await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });

    await setFixture(page);
  },
});

export { expect } from "@playwright/test";
