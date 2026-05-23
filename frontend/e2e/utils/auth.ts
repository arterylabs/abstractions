import { Page } from "@playwright/test";

/**
 * Logs in via the email/password form and waits for the redirect to /events.
 * Call this from tests that need an authenticated session.
 */
export async function login(
  page: Page,
  email: string,
  password: string,
): Promise<void> {
  await page.goto("/login");
  await page.getByTestId("login-email-input").fill(email);
  await page.getByTestId("login-password-input").fill(password);
  await page.getByTestId("login-submit-button").click();
  await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });
}

/**
 * Signs up a new user and waits for the redirect to /dashboard.
 * Use a unique email (e.g. with Date.now()) to avoid conflicts between runs.
 */
export async function signup(
  page: Page,
  name: string,
  email: string,
  password: string,
): Promise<void> {
  await page.goto("/signup");
  await page.getByTestId("signup-name-input").fill(name);
  await page.getByTestId("signup-email-input").fill(email);
  await page.getByTestId("signup-password-input").fill(password);
  await page.getByTestId("signup-submit-button").click();
  await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });
}
