import { Page, expect } from "@playwright/test";

/** Final landing route after signup/login (/dashboard server-redirects here). */
const AUTHED_URL = /\/repositories/;
const AUTH_TIMEOUT_MS = 30_000;

export const apiURL =
  process.env.PLAYWRIGHT_API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:1234";

/**
 * Unique email for a single test run. Avoids collisions across parallel workers.
 */
export function uniqueTestEmail(prefix = "e2e"): string {
  const suffix = `${Date.now()}-${Math.floor(Math.random() * 100000)}`;
  return `${prefix}-${suffix}@test.arterylabs.com`;
}

/** Mark a single test account as verified via the dev-only endpoint. */
export async function verifyTestUser(email: string): Promise<void> {
  const res = await fetch(`${apiURL}/test/verify_user`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(
      `verifyTestUser failed for ${email}: ${res.status} ${text}`,
    );
  }
}

/**
 * Fill a controlled input reliably across browsers (especially WebKit).
 * Asserts the value stuck so React onChange handlers actually fired.
 * Retries: WebKit occasionally accepts fill() without updating the controlled value.
 */
export async function fillField(
  page: Page,
  testId: string,
  value: string,
): Promise<void> {
  const field = page.getByTestId(testId);
  await field.waitFor({ state: "visible" });

  for (let attempt = 0; attempt < 3; attempt++) {
    await field.click();
    await field.fill("");
    await field.fill(value);
    try {
      await expect(field).toHaveValue(value, { timeout: 2_000 });
      return;
    } catch {
      // Value did not stick — retry (common on WebKit under load).
    }
  }

  // Final attempt with character-by-character input, which is more reliable
  // for React controlled inputs when fill() is dropped.
  await field.click();
  await field.fill("");
  await field.pressSequentially(value, { delay: 15 });
  await expect(field).toHaveValue(value);
}

/**
 * Logs in via the email/password form and waits for the authenticated shell.
 */
export async function login(
  page: Page,
  email: string,
  password: string,
): Promise<void> {
  await page.goto("/login");
  await fillField(page, "login-email-input", email);
  await fillField(page, "login-password-input", password);

  await Promise.all([
    page.waitForResponse(
      (r) =>
        r.url().includes("/auth/login") &&
        r.request().method() === "POST" &&
        r.ok(),
      { timeout: AUTH_TIMEOUT_MS },
    ),
    page.getByTestId("login-submit-button").click(),
  ]);

  await settleAuthenticatedSession(page);
}

/**
 * Signs up a new user and waits for the authenticated shell.
 * Use a unique email (e.g. uniqueTestEmail()) to avoid conflicts between runs.
 */
export async function signup(
  page: Page,
  name: string,
  email: string,
  password: string,
): Promise<void> {
  await page.goto("/signup");
  await fillField(page, "signup-name-input", name);
  await fillField(page, "signup-email-input", email);
  await fillField(page, "signup-password-input", password);

  await Promise.all([
    page.waitForResponse(
      (r) =>
        r.url().includes("/auth/signup") &&
        r.request().method() === "POST" &&
        r.ok(),
      { timeout: AUTH_TIMEOUT_MS },
    ),
    page.getByTestId("signup-submit-button").click(),
  ]);

  await settleAuthenticatedSession(page);
}

/** Ensure the auth token cookie was written after login/signup. */
async function assertAuthCookie(page: Page): Promise<string> {
  // Cookie is set asynchronously after the auth response; poll briefly.
  for (let i = 0; i < 20; i++) {
    const cookies = await page.context().cookies();
    const token = cookies.find((c) => c.name === "token");
    if (token?.value) {
      return token.value;
    }
    await page.waitForTimeout(100);
  }
  throw new Error("Auth token cookie was not set after login/signup");
}

/**
 * The API wraps signup/login in a DB transaction that commits after the HTTP
 * response is written. Poll /users/me until the session is actually readable
 * so the SPA doesn't bounce us to /login on the first syncCurrentUser call.
 */
async function waitForSessionReadable(
  page: Page,
  token: string,
): Promise<void> {
  const deadline = Date.now() + AUTH_TIMEOUT_MS;
  let lastStatus = 0;

  while (Date.now() < deadline) {
    const res = await page.request.get(`${apiURL}/users/me`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    lastStatus = res.status();
    if (res.ok()) {
      return;
    }
    await page.waitForTimeout(150);
  }

  throw new Error(
    `Session never became readable via /users/me (last status: ${lastStatus})`,
  );
}

/**
 * Navigate to /repositories, tolerating WebKit races where a concurrent SPA
 * redirect (e.g. to /dashboard) interrupts page.goto.
 */
async function gotoRepositoriesResilient(page: Page): Promise<void> {
  for (let attempt = 0; attempt < 3; attempt++) {
    try {
      await page.goto("/repositories", { waitUntil: "domcontentloaded" });
      return;
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      if (!message.includes("interrupted by another navigation")) {
        throw err;
      }
      // Competing SPA navigation won — wait for it to settle, then decide.
      await page.waitForLoadState("domcontentloaded").catch(() => {});
      if (AUTHED_URL.test(page.url())) {
        return;
      }
      // On /dashboard the server/client will redirect to /repositories.
      if (page.url().includes("/dashboard")) {
        try {
          await page.waitForURL(AUTHED_URL, { timeout: 5_000 });
          return;
        } catch {
          // Retry goto below.
        }
      }
    }
  }
  await page.goto("/repositories", { waitUntil: "domcontentloaded" });
}

/**
 * After login/signup the app navigates to /dashboard → /repositories. Wait
 * until the token cookie is set, the session is readable, and the shell is up.
 *
 * Prefer waiting for the SPA redirect chain over forcing page.goto — WebKit
 * fails hard when goto races a concurrent client navigation
 * ("Navigation ... is interrupted by another navigation").
 */
async function settleAuthenticatedSession(page: Page): Promise<void> {
  const token = await assertAuthCookie(page);
  await waitForSessionReadable(page, token);

  // Give the SPA a chance to leave the auth pages on its own first.
  try {
    await page.waitForURL(
      (url) => {
        const path = url.pathname;
        return !path.includes("/login") && !path.includes("/signup");
      },
      { timeout: 10_000 },
    );
  } catch {
    // Still on login/signup — force navigation below.
  }

  // Wait for the final authenticated landing route.
  try {
    await page.waitForURL(AUTHED_URL, { timeout: 10_000 });
  } catch {
    // SPA never reached /repositories (bounced to login, stuck, etc.).
    await gotoRepositoriesResilient(page);
  }

  // If we landed on /dashboard, wait one more beat for the redirect.
  if (!AUTHED_URL.test(page.url()) && page.url().includes("/dashboard")) {
    try {
      await page.waitForURL(AUTHED_URL, { timeout: 5_000 });
    } catch {
      await gotoRepositoriesResilient(page);
    }
  }

  await waitForAuthenticatedShell(page);
}

/**
 * Wait until the authenticated shell has finished loading (user + workspaces).
 * Prefer this over waitForURL alone — URL can match while client-side auth still
 * fails and redirects to /login.
 */
export async function waitForAuthenticatedShell(page: Page): Promise<void> {
  // Desktop: sidebar nav. Mobile: hamburger trigger (nav is inside a closed sheet).
  const shell = page
    .getByTestId("dashboard-nav")
    .or(page.getByRole("button", { name: /toggle sidebar/i }));

  await expect(shell.first()).toBeVisible({ timeout: AUTH_TIMEOUT_MS });
  await expect(page).not.toHaveURL(/\/login/);
}

/**
 * Creates a new workspace via the switcher's "create workspace" flow and
 * waits for the redirect back to the authenticated shell.
 */
export async function createWorkspace(page: Page, name: string): Promise<void> {
  await page.getByTestId("workspace-switcher-trigger").click();
  await page.getByTestId("create-workspace-link").click();
  await page.getByTestId("create-workspace-name-input").fill(name);
  await page.getByTestId("create-workspace-submit").click();
  await page.waitForURL(/\/(repositories|dashboard)/, { timeout: 15_000 });
}
