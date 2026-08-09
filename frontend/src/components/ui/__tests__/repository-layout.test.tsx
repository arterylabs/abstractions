import RepositoryPage from "@/app/(dashboard)/repositories/[id]/page";
import { getRepositorySectionValue } from "@/app/(dashboard)/repositories/[id]/utils";
import { getNewsletterSectionValue } from "@/app/(dashboard)/repositories/[id]/newsletters/utils";

const mockRedirect = jest.fn();

jest.mock("next/navigation", () => ({
  redirect: (...args: unknown[]) => mockRedirect(...args),
}));

describe("repository layout routing", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("maps overview routes to the overview tab", () => {
    expect(
      getRepositorySectionValue("/repositories/repo-1/overview", "repo-1"),
    ).toBe("overview");
  });

  it("maps nested newsletter routes to the newsletters tab", () => {
    expect(
      getRepositorySectionValue(
        "/repositories/repo-1/newsletters/newsletter-1",
        "repo-1",
      ),
    ).toBe("newsletters");
  });

  it("maps settings routes to the settings tab", () => {
    expect(
      getRepositorySectionValue("/repositories/repo-1/settings", "repo-1"),
    ).toBe("settings");
  });

  it("defaults to overview for unknown paths", () => {
    expect(
      getRepositorySectionValue("/repositories/repo-1/unknown", "repo-1"),
    ).toBe("overview");
  });

  it("redirects the repository root page to overview", async () => {
    await RepositoryPage({
      params: Promise.resolve({ id: "repo-1" }),
    });

    expect(mockRedirect).toHaveBeenCalledWith("/repositories/repo-1/overview");
  });
});

describe("newsletters sub-tab routing", () => {
  it("maps the newsletters root to issues", () => {
    expect(
      getNewsletterSectionValue("/repositories/repo-1/newsletters", "repo-1"),
    ).toBe("issues");
  });

  it("maps an issue detail route to issues", () => {
    expect(
      getNewsletterSectionValue(
        "/repositories/repo-1/newsletters/newsletter-1",
        "repo-1",
      ),
    ).toBe("issues");
  });

  it("maps the topics route to topics", () => {
    expect(
      getNewsletterSectionValue(
        "/repositories/repo-1/newsletters/topics",
        "repo-1",
      ),
    ).toBe("topics");
  });

  it("maps the schedule route to schedule", () => {
    expect(
      getNewsletterSectionValue(
        "/repositories/repo-1/newsletters/schedule",
        "repo-1",
      ),
    ).toBe("schedule");
  });

  it("maps the recipients route to recipients", () => {
    expect(
      getNewsletterSectionValue(
        "/repositories/repo-1/newsletters/recipients",
        "repo-1",
      ),
    ).toBe("recipients");
  });
});
