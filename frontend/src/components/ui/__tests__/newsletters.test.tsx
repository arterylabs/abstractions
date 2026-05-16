import { render, screen } from "@testing-library/react";
import RepositoryNewslettersPage from "@/app/(dashboard)/repositories/[id]/newsletters/page";
import { NewsletterPreviewFrame } from "@/app/(dashboard)/repositories/[id]/newsletters/[newsletterId]/page";
import { NewsletterCard } from "@/components/newsletters/NewsletterCard";
import { NextUpCard } from "@/components/newsletters/NextUpCard";
import { formatRate } from "@/lib/format";

const mockPush = jest.fn();
const mockList = jest.fn();
const mockTopicsList = jest.fn();
const mockGet = jest.fn();
const mockDispatch = jest.fn();

let mockParams = { id: "repo-1", newsletterId: "newsletter-1" };

jest.mock("@/store/hooks", () => ({
  useAppDispatch: () => mockDispatch,
  useAppSelector: (selector: (state: unknown) => unknown) =>
    selector({
      auth: { currentWorkspace: { id: "workspace-1" } },
      repositories: { items: [] },
    }),
}));

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
  useParams: () => mockParams,
}));

jest.mock("@/api", () => ({
  __esModule: true,
  default: {
    newsletters: {
      list: (...args: unknown[]) => mockList(...args),
      topics: {
        list: (...args: unknown[]) => mockTopicsList(...args),
      },
      get: (...args: unknown[]) => mockGet(...args),
    },
  },
}));

describe("Repository newsletters", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockParams = { id: "repo-1", newsletterId: "newsletter-1" };
    mockDispatch.mockImplementation((action: unknown) => action);
  });

  it("formats rates as percentages", () => {
    expect(formatRate(0.4)).toBe("40.0%");
  });

  it("displays newsletter engagement counts in the newsletter list", async () => {
    mockList.mockResolvedValue({
      data: [
        {
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Weekly roundup",
          sentAt: "2025-01-01T00:00:00.000Z",
          recipientCount: 10,
          openCount: 4,
          clickCount: 2,
          openRate: 0.4,
          clickRate: 0.2,
          createdAt: "2025-01-01T00:00:00.000Z",
        },
      ],
    });
    mockTopicsList.mockResolvedValue({
      data: [
        {
          id: "topic-1",
          repoId: "repo-1",
          title: "Topic A",
          description: "Queue topic",
          relatedFiles: [],
          sentAt: null,
          createdAt: "2025-01-02T00:00:00.000Z",
        },
      ],
    });

    render(<RepositoryNewslettersPage />);

    expect(await screen.findByText("Weekly roundup")).toBeInTheDocument();
    expect(screen.getByText("4 opens")).toBeInTheDocument();
    expect(screen.getByText("2 clicks")).toBeInTheDocument();
  });

  it("renders the newsletter preview iframe with a sandbox attribute", () => {
    render(<NewsletterPreviewFrame htmlContent="<p>Hello world</p>" />);

    expect(screen.getByTitle("Newsletter Preview")).toHaveAttribute(
      "sandbox",
      "allow-same-origin",
    );
  });

  it("renders queued newsletter topics", async () => {
    mockList.mockResolvedValue({ data: [] });
    mockTopicsList.mockResolvedValue({
      data: [
        {
          id: "topic-1",
          repoId: "repo-1",
          title: "Topic A",
          description: "Queue topic",
          relatedFiles: [],
          sentAt: null,
          createdAt: "2025-01-02T00:00:00.000Z",
        },
      ],
    });

    render(<RepositoryNewslettersPage />);

    expect(await screen.findByText("Topic A")).toBeInTheDocument();
    expect(screen.getByText("Pending")).toBeInTheDocument();
  });

  it("renders a newsletter card summary", () => {
    render(
      <NewsletterCard
        newsletter={{
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Weekly roundup",
          sentAt: "2025-01-01T00:00:00.000Z",
          recipientCount: 10,
          openCount: 4,
          clickCount: 2,
          openRate: 0.4,
          clickRate: 0.2,
          createdAt: "2025-01-01T00:00:00.000Z",
        }}
        onClick={jest.fn()}
      />,
    );

    expect(
      screen.getByRole("button", { name: /weekly roundup/i }),
    ).toBeInTheDocument();
    expect(screen.getByText("4 opens")).toBeInTheDocument();
    expect(screen.getByText("2 clicks")).toBeInTheDocument();
  });

  it("renders a next up card for a queued topic", () => {
    render(
      <NextUpCard
        nextTopic={{
          id: "topic-1",
          repoId: "repo-1",
          title: "Topic A",
          description: "Queue topic",
          relatedFiles: [],
          position: 0,
          sentAt: null,
          requestedByUserId: null,
          requestedByName: null,
          requestedByAvatar: null,
          newsletterId: null,
          createdAt: "2025-01-02T00:00:00.000Z",
        }}
        repository={{
          id: "repo-1",
          workspaceId: "workspace-1",
          userId: "user-1",
          fullName: "acme/repo",
          description: "",
          defaultBranch: "main",
          status: "indexed",
          indexedAt: null,
          sendDay: 1,
          sendHour: 9,
          newsletterEnabled: true,
          lastNewsletterSentAt: null,
          emails: [],
          sendToWorkspaceMembers: true,
          createdAt: "2025-01-01T00:00:00.000Z",
          updatedAt: "2025-01-01T00:00:00.000Z",
        }}
        isLoading={false}
        onPreview={jest.fn()}
        onSendNow={jest.fn()}
      />,
    );

    expect(screen.getByText("Topic A")).toBeInTheDocument();
    expect(screen.getByText("Queue topic")).toBeInTheDocument();
    expect(screen.getByText(/Scheduled to send/)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Preview" })).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Send Now" }),
    ).toBeInTheDocument();
  });
});
