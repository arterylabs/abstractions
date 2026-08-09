import { render, screen, waitFor } from "@testing-library/react";
import RepositoriesPage from "@/app/(dashboard)/repositories/page";
import { RepositoryStatusBadge } from "@/components/repositories/RepositoryStatusBadge";

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

  it("renders an indexed badge with the success variant", () => {
    render(<RepositoryStatusBadge status="indexed" />);

    expect(screen.getByText("Indexed")).toHaveAttribute(
      "data-variant",
      "success",
    );
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
        status: "indexing",
        indexedAt: null,
        sendDay: 1,
        sendHour: 9,
        newsletterSentAt: null,
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
      screen.getAllByRole("button", { name: /connect repository/i })[0],
    ).toBeInTheDocument();
    expect(
      screen.getByText("Connect your first repository"),
    ).toBeInTheDocument();
  });
});
