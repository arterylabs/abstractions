import { render, waitFor } from "@testing-library/react";
import RepositoriesPage from "@/app/(dashboard)/repositories/page";

const mockDispatch = jest.fn();
const mockPush = jest.fn();

type MockState = {
  auth: {
    currentWorkspace: {
      id: string;
    };
  };
  repositories: {
    items: Array<Record<string, unknown>>;
    status: "idle" | "loading" | "success" | "failed";
    error: string | null;
  };
};

let mockState: MockState = {
  auth: {
    currentWorkspace: {
      id: "workspace-1",
    },
  },
  repositories: {
    items: [],
    status: "idle",
    error: null,
  },
};

jest.mock("@/store/hooks", () => ({
  useAppDispatch: () => mockDispatch,
  useAppSelector: (selector: (state: MockState) => unknown) =>
    selector(mockState),
}));

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

describe("RepositoriesPage indexing poll", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockState = {
      auth: {
        currentWorkspace: {
          id: "workspace-1",
        },
      },
      repositories: {
        items: [
          {
            id: "repo-1",
            workspaceId: "workspace-1",
            userId: "user-1",
            githubRepoId: 101,
            fullName: "acme/indexing-repo",
            defaultBranch: "main",
            status: "indexing",
            indexedAt: null,
            sendDay: 1,
            sendHour: 9,
            newsletterSentAt: null,
            createdAt: "2025-01-01T00:00:00.000Z",
            updatedAt: "2025-01-01T00:00:00.000Z",
          },
        ],
        status: "success",
        error: null,
      },
    };
  });

  it("sets a polling interval when any repo is indexing", async () => {
    const setIntervalSpy = jest.spyOn(window, "setInterval");

    render(<RepositoriesPage />);

    await waitFor(() => {
      expect(setIntervalSpy).toHaveBeenCalledWith(expect.any(Function), 10000);
    });

    setIntervalSpy.mockRestore();
  });
});
