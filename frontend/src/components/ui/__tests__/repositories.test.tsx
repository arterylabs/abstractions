import { render, screen, waitFor } from "@testing-library/react";
import RepositoriesPage, {
  RepositoryStatusBadge,
} from "@/app/(dashboard)/repositories/page";

const mockDispatch = jest.fn();

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

describe("RepositoriesPage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockState = {
      auth: {
        currentWorkspace: {
          id: "workspace-1",
        },
      },
      repositories: {
        items: [],
        status: "success",
        error: null,
      },
    };
  });

  it("renders an indexing badge with animate-pulse", () => {
    render(<RepositoryStatusBadge status="indexing" />);

    expect(screen.getByText("Indexing…")).toHaveClass("animate-pulse");
  });

  it("renders a ready badge with the success variant", () => {
    render(<RepositoryStatusBadge status="indexed" />);

    expect(screen.getByText("Ready")).toHaveClass("bg-success");
  });

  it("sets a polling interval when any repo is indexing", async () => {
    const setIntervalSpy = jest.spyOn(window, "setInterval");

    mockState.repositories.items = [
      {
        id: "repo-1",
        workspaceId: "workspace-1",
        userId: "user-1",
        githubRepoId: 101,
        fullName: "acme/indexing-repo",
        defaultBranch: "main",
        webhookId: null,
        status: "indexing",
        indexedAt: null,
        sendDay: 1,
        sendHour: 9,
        lastNewsletterSentAt: null,
        createdAt: "2025-01-01T00:00:00.000Z",
        updatedAt: "2025-01-01T00:00:00.000Z",
      },
    ];

    render(<RepositoriesPage />);

    await waitFor(() => {
      expect(setIntervalSpy).toHaveBeenCalledWith(expect.any(Function), 10000);
    });

    setIntervalSpy.mockRestore();
  });

  it("renders the empty state CTA button", () => {
    render(<RepositoriesPage />);

    expect(
      screen.getAllByRole("link", { name: "Connect Repository" })[0],
    ).toBeInTheDocument();
    expect(screen.getByText("No repositories connected")).toBeInTheDocument();
  });
});
