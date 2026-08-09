import { render, waitFor } from "@testing-library/react";
import { RepositoryIndexingStatusBanner } from "@/components/repositories/RepositoryIndexingStatusBanner";

const mockDispatch = jest.fn();
const mockLoadStatus = jest.fn();

let mockStatusResult: {
  result?: {
    status: "pending" | "indexing" | "indexed" | "error";
    totalFiles: number;
    summarizedFiles: number;
    progressPercentage: number;
    languageBreakdown: Array<{
      language: string;
      fileCount: number;
      percentage: number;
    }>;
    indexedAt?: string | null;
    indexError?: string | null;
  };
  isFetching: boolean;
} = {
  result: undefined,
  isFetching: false,
};

jest.mock("@/store/hooks", () => ({
  useAppDispatch: () => mockDispatch,
}));

jest.mock("@/api/hooks", () => ({
  __esModule: true,
  default: () => [mockStatusResult, mockLoadStatus],
}));

jest.mock("@/api", () => ({
  __esModule: true,
  default: {
    repositories: {
      getStatus: jest.fn(),
    },
  },
}));

describe("RepositoryIndexingStatusBanner", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("polls the status endpoint and refreshes the repository when indexed", async () => {
    const setIntervalSpy = jest.spyOn(window, "setInterval");
    mockStatusResult = {
      result: {
        status: "indexed",
        totalFiles: 12,
        summarizedFiles: 12,
        progressPercentage: 100,
        languageBreakdown: [],
      },
      isFetching: false,
    };

    render(
      <RepositoryIndexingStatusBanner
        workspaceId="workspace-1"
        repository={{
          id: "repo-1",
          workspaceId: "workspace-1",
          userId: "user-1",
          fullName: "acme/repo",
          description: "",
          defaultBranch: "main",
          status: "indexing",
          indexedAt: null,
          indexError: null,
          sendDay: 1,
          sendHour: 9,
          newsletterEnabled: true,
          lastNewsletterSentAt: null,
          emails: [],
          sendToWorkspaceMembers: true,
          createdAt: "2025-01-01T00:00:00Z",
          updatedAt: "2025-01-01T00:00:00Z",
        }}
      />,
    );

    await waitFor(() => {
      expect(mockLoadStatus).toHaveBeenCalledWith("workspace-1", "repo-1");
    });

    expect(setIntervalSpy).toHaveBeenCalledWith(expect.any(Function), 5000);
    expect(mockDispatch).toHaveBeenCalled();

    setIntervalSpy.mockRestore();
  });
});
