import { act, render, screen } from "@testing-library/react";
import RepositoryNewslettersPage from "@/app/(dashboard)/repositories/[id]/newsletters/page";
import { NewsletterPreviewFrame } from "@/app/(dashboard)/repositories/[id]/newsletters/[newsletterId]/NewsletterPreviewFrame";
import { NewsletterCard } from "@/components/newsletters/NewsletterCard";

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
      resendNewsletter: jest.fn(),
    },
  },
}));

describe("Repository newsletters", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockParams = { id: "repo-1", newsletterId: "newsletter-1" };
    mockDispatch.mockImplementation((action: unknown) => action);
  });

  it("displays newsletter subject in the newsletter list", async () => {
    mockList.mockResolvedValue({
      data: [
        {
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Weekly roundup",
          status: "sent",
          errorMessage: null,
          sentAt: "2025-01-01T00:00:00.000Z",
          recipientCount: 10,
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
  });

  it("renders the newsletter preview iframe with a sandbox attribute", () => {
    render(<NewsletterPreviewFrame htmlContent="<p>Hello world</p>" />);

    expect(screen.getByTitle("Newsletter Preview")).toHaveAttribute(
      "sandbox",
      "allow-same-origin",
    );
  });

  // The archive summarises the queue rather than listing it — topics are
  // managed one click away in the Newsletters sidebar.
  it("summarises the queued topics", async () => {
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

    expect(
      await screen.findByText("1 topic queued for the next issue"),
    ).toBeInTheDocument();
  });

  it("renders a newsletter card summary", () => {
    render(
      <NewsletterCard
        newsletter={{
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Weekly roundup",
          status: "sent",
          errorMessage: null,
          sentAt: "2025-01-01T00:00:00.000Z",
          recipientCount: 10,
          createdAt: "2025-01-01T00:00:00.000Z",
        }}
        onClick={jest.fn()}
      />,
    );

    expect(
      screen.getByRole("button", { name: /weekly roundup/i }),
    ).toBeInTheDocument();
  });

  it("renders a status badge on the newsletter card", () => {
    render(
      <NewsletterCard
        newsletter={{
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Weekly roundup",
          status: "error",
          errorMessage: "AI provider unavailable",
          sentAt: null,
          recipientCount: 0,
          createdAt: "2025-01-01T00:00:00.000Z",
        }}
        onClick={jest.fn()}
        onResend={jest.fn()}
      />,
    );

    expect(screen.getByText("Error")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "More actions" }),
    ).toBeInTheDocument();
  });
});

describe("Newsletter status polling", () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.clearAllMocks();
    mockParams = { id: "repo-1", newsletterId: "newsletter-1" };
    mockDispatch.mockImplementation((action: unknown) => action);
    mockTopicsList.mockResolvedValue({ data: [] });
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("polls for newsletter status while a non-terminal newsletter exists", async () => {
    // First call returns a "queued" newsletter; subsequent calls return "sent"
    mockList
      .mockResolvedValueOnce({
        data: [
          {
            id: "newsletter-1",
            repoId: "repo-1",
            subject: "In progress",
            status: "queued",
            errorMessage: null,
            sentAt: null,
            recipientCount: 0,
            createdAt: "2025-01-01T00:00:00.000Z",
          },
        ],
      })
      .mockResolvedValue({
        data: [
          {
            id: "newsletter-1",
            repoId: "repo-1",
            subject: "In progress",
            status: "sent",
            errorMessage: null,
            sentAt: "2025-01-01T01:00:00.000Z",
            recipientCount: 5,
            createdAt: "2025-01-01T00:00:00.000Z",
          },
        ],
      });

    render(<RepositoryNewslettersPage />);

    // Wait for the initial fetch to resolve
    await act(async () => {
      await Promise.resolve();
    });

    const callsAfterInitial = mockList.mock.calls.length;

    // Advance past one poll interval
    await act(async () => {
      jest.advanceTimersByTime(5000);
      await Promise.resolve();
    });

    expect(mockList.mock.calls.length).toBeGreaterThan(callsAfterInitial);
  });

  it("does not poll when all newsletters are in a terminal status", async () => {
    mockList.mockResolvedValue({
      data: [
        {
          id: "newsletter-1",
          repoId: "repo-1",
          subject: "Already sent",
          status: "sent",
          errorMessage: null,
          sentAt: "2025-01-01T01:00:00.000Z",
          recipientCount: 5,
          createdAt: "2025-01-01T00:00:00.000Z",
        },
      ],
    });

    render(<RepositoryNewslettersPage />);

    await act(async () => {
      await Promise.resolve();
    });

    const callsAfterInitial = mockList.mock.calls.length;

    // Advance well past the poll interval — no additional calls expected
    await act(async () => {
      jest.advanceTimersByTime(15000);
      await Promise.resolve();
    });

    expect(mockList.mock.calls.length).toBe(callsAfterInitial);
  });
});
