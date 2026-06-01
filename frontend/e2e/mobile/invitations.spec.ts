import { test, expect } from "../fixtures";

test.describe("Mobile invitations", () => {
  test("invitations page renders send button below search on mobile", async ({
    authedPage,
  }) => {
    await authedPage.goto("/settings/workspace/invitations");

    const sendButton = authedPage.getByTestId("send-invitations-button");
    const searchBar = authedPage.getByPlaceholder("Search invitations...");

    await expect(sendButton).toBeVisible();
    await expect(searchBar).toBeVisible();

    // On mobile the button stacks below the search bar
    const buttonBox = await sendButton.boundingBox();
    const searchBox = await searchBar.boundingBox();

    expect(buttonBox).not.toBeNull();
    expect(searchBox).not.toBeNull();

    // button should be below (or at the same top) as the search bar on mobile
    expect(buttonBox!.y).toBeGreaterThanOrEqual(searchBox!.y);
  });

  test("can open invite dialog on mobile", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await authedPage.getByTestId("send-invitations-button").click();

    await expect(
      authedPage.getByText("Send Workspace Invitations"),
    ).toBeVisible();
  });

  test("invite dialog accepts email on mobile", async ({ authedPage }) => {
    await authedPage.goto("/settings/workspace/invitations");

    await authedPage.getByTestId("send-invitations-button").click();

    const inviteeEmail = `mobile-invitee-${Date.now()}@test.arterylabs.com`;
    await authedPage.getByTestId("invite-emails-textarea").fill(inviteeEmail);

    await authedPage.getByTestId("send-invitations-dialog-button").click();

    await expect(
      authedPage.getByText("Invitations sent successfully", { exact: false }),
    ).toBeVisible({ timeout: 10_000 });
  });
});
