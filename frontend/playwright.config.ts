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

  // Cap parallelism: each authed test signs up a new user + workspace, and the
  // API + local DB get flaky under high concurrent signup load (especially WebKit).
  workers: process.env.CI ? 1 : 3,

  // Authed fixtures sign up + verify + reload; give them room on slower browsers.
  timeout: 60_000,
  expect: {
    timeout: 10_000,
  },

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
      // Mobile viewport specs belong only in the mobile projects below.
      testIgnore: "**/mobile/**/*.spec.ts",
    },
    {
      name: "mobile-chrome",
      use: { ...devices["Pixel 5"] },
      testMatch: "**/mobile/**/*.spec.ts",
    },
    {
      name: "mobile-safari",
      use: { ...devices["iPhone 13"] },
      testMatch: "**/mobile/**/*.spec.ts",
    },
    // Uncomment to run on additional browsers:
    // { name: "firefox", use: { ...devices["Desktop Firefox"] } },
    // { name: "webkit",  use: { ...devices["Desktop Safari"] } },
  ],
});
