import { test, expect } from "../fixtures";
import {
  signup,
  uniqueTestEmail,
  verifyTestUser,
  waitForAuthenticatedShell,
} from "../utils/auth";

test.describe("Workspace invitations", () => {
  test("owner can send an invitation", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await authedPage.getByTestId("send-invitations-button").click();
    await expect(
      authedPage.getByText("Send Workspace Invitations"),
    ).toBeVisible();

    const inviteeEmail = uniqueTestEmail("invitee");
    await authedPage.getByTestId("invite-emails-textarea").fill(inviteeEmail);
    await authedPage.getByTestId("send-invitations-dialog-button").click();

    await expect(
      authedPage.getByText("Invitations sent successfully", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });
  });

  test("invite dialog shows error for invalid email", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await authedPage.getByTestId("send-invitations-button").click();
    await authedPage.getByTestId("invite-emails-textarea").fill("not-an-email");
    await authedPage.getByTestId("send-invitations-dialog-button").click();

    await expect(
      authedPage.getByText("Invalid email format", { exact: false }),
    ).toBeVisible();
  });

  test("invitee can accept an invitation and join the workspace", async ({
    browser,
    authedPage,
  }) => {
    const workspaceName = await authedPage
      .getByTestId("workspace-switcher-name")
      .textContent();

    const inviteeEmail = uniqueTestEmail("accept");
    const inviteePassword = "Password123!";

    await authedPage.goto("/settings/workspace/invitations");
    await authedPage.getByTestId("send-invitations-button").click();
    await authedPage.getByTestId("invite-emails-textarea").fill(inviteeEmail);
    await authedPage.getByTestId("send-invitations-dialog-button").click();
    await expect(
      authedPage.getByText("Invitations sent successfully", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });

    // Invitee signs up in a separate browser context (separate cookies).
    const inviteeContext = await browser.newContext();
    const inviteePage = await inviteeContext.newPage();

    await signup(inviteePage, "Invitee User", inviteeEmail, inviteePassword);
    await verifyTestUser(inviteeEmail);
    await inviteePage.goto("/me/invitations");
    await waitForAuthenticatedShell(inviteePage).catch(() => {
      // /me/invitations may not show dashboard-nav; ignore shell wait failure.
    });

    // Pending invitations page (may not mount dashboard shell).
    await expect(
      inviteePage.getByRole("heading", { name: /Pending invitations/i }),
    ).toBeVisible({ timeout: 15_000 });

    const acceptButton = inviteePage.getByRole("button", {
      name: "Accept Invitation",
    });
    await expect(acceptButton).toBeVisible({ timeout: 15_000 });
    await acceptButton.click();

    await expect(
      inviteePage.getByText(/Invitation Accepted|added to the workspace/i),
    ).toBeVisible({ timeout: 15_000 });

    await inviteePage.waitForURL(/\/(repositories|dashboard)/, {
      timeout: 15_000,
    });
    await inviteePage.goto("/repositories");
    await expect(
      inviteePage.getByTestId("workspace-switcher-trigger"),
    ).toContainText(workspaceName?.trim() ?? "", { timeout: 15_000 });

    await inviteeContext.close();
  });
});
