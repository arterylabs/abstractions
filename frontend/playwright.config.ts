import { defineConfig, devices } from "@playwright/test";

/**
 * Base URL is read from the environment so the same config works locally
 * (http://localhost:3000) and in CI (pointing at a deployed preview or the
 * locally-started server).
 */
const baseURL = process.env.PLAYWRIGHT_BASE_URL || "http://localhost:3000";

export default defineConfig({
  testDir: "./e2e",
  globalTeardown: "./e2e/global-teardown.ts",

  // Run all tests in parallel within each file.
  fullyParallel: true,

  // Fail the build if test.only is accidentally committed.
  forbidOnly: !!process.env.CI,

  // Retry failing tests twice in CI; no retries locally.
  retries: process.env.CI ? 2 : 0,

  // Use a single worker in CI to avoid shared-state flakiness across parallel
  // browser instances; locally use the default (all CPU cores).
  workers: process.env.CI ? 1 : undefined,

  reporter: [
    ["list"],
    // HTML report is written to playwright-report/ and uploaded as a CI artifact.
    ["html", { outputFolder: "playwright-report", open: "never" }],
  ],

  use: {
    baseURL,

    // Capture traces and videos only on the first retry so artifacts don't
    // balloon for consistently-passing tests.
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    video: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    // Uncomment to run on additional browsers:
    // { name: "firefox", use: { ...devices["Desktop Firefox"] } },
    // { name: "webkit",  use: { ...devices["Desktop Safari"] } },
  ],
});
