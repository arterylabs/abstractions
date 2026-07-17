import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ConnectPage from "@/app/(dashboard)/repositories/connect/page";

const mockPush = jest.fn();
const mockDispatch = jest.fn();
const mockListRepos = jest.fn();
const mockConnectRepository = jest.fn();
const mockGetWorkspaceHealth = jest.fn();

let mockState = {
  auth: {
    currentWorkspace: {
      id: "workspace-1",
    },
  },
};

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

jest.mock("@/store/hooks", () => ({
  useAppDispatch: () => mockDispatch,
  useAppSelector: (selector: (state: typeof mockState) => unknown) =>
    selector(mockState),
}));

jest.mock("@/api", () => ({
  __esModule: true,
  default: {
    github: {
      listRepos: (...args: unknown[]) => mockListRepos(...args),
      connectRepository: (...args: unknown[]) => mockConnectRepository(...args),
    },
    workspaces: {
      getWorkspaceHealth: (...args: unknown[]) =>
        mockGetWorkspaceHealth(...args),
    },
  },
}));

describe("ConnectPage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockGetWorkspaceHealth.mockResolvedValue({ data: null });
    mockState = {
      auth: {
        currentWorkspace: {
          id: "workspace-1",
        },
      },
    };
  });

  it("renders a loading spinner while repos are fetching", async () => {
    mockListRepos.mockImplementation(() => new Promise(() => undefined));

    render(<ConnectPage />);

    expect(
      await screen.findByText("Loading repositories..."),
    ).toBeInTheDocument();
  });

  it("renders the repo list after data loads", async () => {
    mockListRepos.mockResolvedValue({
      data: [
        {
          id: 1,
          fullName: "acme/repo-one",
          defaultBranch: "main",
          private: true,
          description: "First repository",
          htmlUrl: "https://github.com/acme/repo-one",
        },
      ],
    });

    render(<ConnectPage />);

    expect(await screen.findByText("acme/repo-one")).toBeInTheDocument();
    expect(screen.getByText("First repository")).toBeInTheDocument();
    expect(screen.getByText("Private")).toBeInTheDocument();
  });

  it("calls connectRepository when clicking Connect", async () => {
    const user = userEvent.setup();

    mockListRepos.mockResolvedValue({
      data: [
        {
          id: 42,
          fullName: "acme/repo-two",
          defaultBranch: "develop",
          private: false,
          description: "Second repository",
          htmlUrl: "https://github.com/acme/repo-two",
        },
      ],
    });
    mockConnectRepository.mockResolvedValue({ data: {} });

    render(<ConnectPage />);

    await user.click(await screen.findByRole("button", { name: "Connect" }));

    await waitFor(() => {
      expect(mockConnectRepository).toHaveBeenCalledWith("workspace-1", {
        githubRepoId: 42,
        fullName: "acme/repo-two",
        defaultBranch: "develop",
        description: "Second repository",
      });
    });
  });
});
