import { render, screen } from "@testing-library/react";
import RepositoryNewslettersPage, {
  formatRate,
} from "@/app/(dashboard)/repositories/[id]/newsletters/page";
import { NewsletterPreviewFrame } from "@/app/(dashboard)/repositories/[id]/newsletters/[newsletterId]/page";

const mockPush = jest.fn();
const mockList = jest.fn();
const mockTopicsList = jest.fn();
const mockGet = jest.fn();

let mockParams = { id: "repo-1", newsletterId: "newsletter-1" };

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
  });

  it("formats rates as percentages", () => {
    expect(formatRate(0.4)).toBe("40.0%");
  });

  it("displays open and click rates in the newsletter list", async () => {
    mockList.mockResolvedValue({
      data: [
        {
          id: "newsletter-1",
          repoId: "repo-1",
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

    expect(await screen.findByText("40.0%")).toBeInTheDocument();
    expect(screen.getByText("20.0%")).toBeInTheDocument();
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
});
