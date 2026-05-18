import { test, expect } from "../fixtures";

test.describe("Workspace invitations", () => {
  test("invitations page shows send invitations button", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await expect(
      authedPage.getByTestId("send-invitations-button"),
    ).toBeVisible();
  });

  test("owner can open invite dialog and send an invitation", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await authedPage.getByTestId("send-invitations-button").click();

    // Dialog should appear
    await expect(
      authedPage.getByText("Send Workspace Invitations"),
    ).toBeVisible();

    const inviteeEmail = `invitee-${Date.now()}@test.arterylabs.com`;
    await authedPage.getByTestId("invite-emails-textarea").fill(inviteeEmail);

    await authedPage.getByTestId("send-invitations-dialog-button").click();

    // Success toast should appear and dialog should close
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
});
