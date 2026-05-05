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
      email: `e2e-${suffix}@example.com`,
      password: "Password123!",
    });
  },

  authedPage: async ({ page, user }, setFixture) => {
    await signup(page, user.name, user.email, user.password);
    await setFixture(page);
  },
});

export { expect } from "@playwright/test";
