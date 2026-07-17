import { test, expect } from "../fixtures";

test.describe("User account settings", () => {
  test("user can navigate to account settings", async ({ authedPage }) => {
    await authedPage.getByTestId("nav-user-trigger").click();
    await authedPage.getByTestId("account-link").click();

    await expect(authedPage).toHaveURL(/\/account/);
    await expect(
      authedPage.getByRole("heading", { name: "Account Settings" }),
    ).toBeVisible();
  });

  test("account settings page shows name and email fields", async ({
    authedPage,
    user,
  }) => {
    await authedPage.goto("/account");

    await expect(authedPage.getByTestId("account-name-input")).toBeVisible();
    await expect(authedPage.getByTestId("account-name-input")).toHaveValue(
      user.name,
    );

    // Email field exists but is read-only
    const emailInput = authedPage.getByTestId("account-email-input");
    await expect(emailInput).toBeVisible();
    await expect(emailInput).toBeDisabled();
  });

  test("user can update their display name", async ({ authedPage }) => {
    const updatedName = `Updated Name ${Date.now()}`;

    await authedPage.goto("/account");

    const nameInput = authedPage.getByTestId("account-name-input");
    await nameInput.clear();
    await nameInput.fill(updatedName);

    await authedPage.getByTestId("account-save-button").click();

    await expect(
      authedPage.getByText("Profile saved", { exact: false }),
    ).toBeVisible();

    // Reload and verify the name persisted
    await authedPage.reload();
    await expect(authedPage.getByTestId("account-name-input")).toHaveValue(
      updatedName,
    );
  });

  test("account settings shows sidebar navigation tabs", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account");

    await expect(
      authedPage.getByRole("tab", { name: /General/i }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("tab", { name: /Security/i }),
    ).toBeVisible();
    await expect(
      authedPage.getByRole("tab", { name: /Advanced/i }),
    ).toBeVisible();
  });

  test("user can navigate to security tab", async ({ authedPage }) => {
    await authedPage.goto("/account");

    await authedPage.getByRole("tab", { name: /Security/i }).click();

    await expect(authedPage).toHaveURL(/\/account\/security/);
  });

  test("security page shows password management section", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/security");

    // Password card is always present
    await expect(
      authedPage.getByRole("heading", { name: "Password" }),
    ).toBeVisible();

    // For a freshly created account (password sign-up), the button text is either
    // "Change Password" or "Setup Password" depending on the authenticator state.
    const passwordButton = authedPage
      .getByTestId("change-password-button")
      .or(authedPage.getByTestId("setup-password-button"));
    await expect(passwordButton).toBeVisible();
  });

  test("security page shows Google account linking section", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/security");

    await expect(
      authedPage.getByRole("heading", { name: "Google Account" }),
    ).toBeVisible();
  });

  test("can open password change modal from security settings", async ({
    authedPage,
  }) => {
    await authedPage.goto("/account/security");

    const passwordButton = authedPage
      .getByTestId("change-password-button")
      .or(authedPage.getByTestId("setup-password-button"));

    await passwordButton.click();

    await expect(authedPage.getByRole("dialog")).toBeVisible();
    await expect(
      authedPage.getByRole("heading", { name: "Change Password" }),
    ).toBeVisible();
  });

  test("can navigate directly to security settings via switcher", async ({
    authedPage,
  }) => {
    await authedPage.getByTestId("nav-user-trigger").click();
    await authedPage.getByTestId("account-security-link").click();

    await expect(authedPage).toHaveURL(/\/account\/security/);
  });
});
