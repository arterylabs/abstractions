import { act, render } from "@testing-library/react";
import RepositoryNewslettersPage from "@/app/(dashboard)/repositories/[id]/newsletters/page";

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

    await act(async () => {
      await Promise.resolve();
    });

    const callsAfterInitial = mockList.mock.calls.length;

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

    await act(async () => {
      jest.advanceTimersByTime(15000);
      await Promise.resolve();
    });

    expect(mockList.mock.calls.length).toBe(callsAfterInitial);
  });
});
