import { render, screen } from "@testing-library/react";
import Grid from "../grid";

// Virtualizer needs a real scroll container size; mock it so row bodies
// render deterministically in jsdom.
jest.mock("@tanstack/react-virtual", () => ({
  useVirtualizer: ({ count }: { count: number }) => ({
    getTotalSize: () => count * 40,
    getVirtualItems: () =>
      Array.from({ length: count }, (_, index) => ({
        key: index,
        index,
        start: index * 40,
        size: 40,
      })),
  }),
}));

const columns = [{ name: "Name" }, { name: "Status" }];
const rows = [
  ["Alice", "active"],
  ["Bob", "inactive"],
];

describe("Grid", () => {
  it("renders without crashing", () => {
    render(<Grid columns={columns} rows={rows} height={300} />);
  });

  it("renders column headers", () => {
    const { getByText } = render(
      <Grid columns={columns} rows={rows} height={300} />,
    );
    expect(getByText("Name")).toBeInTheDocument();
    expect(getByText("Status")).toBeInTheDocument();
  });

  it("renders row cell values", () => {
    render(<Grid columns={columns} rows={rows} height={300} />);
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("active")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
    expect(screen.getByText("inactive")).toBeInTheDocument();
  });

  it("shows skeleton placeholders when totalRecords exceeds loaded rows", () => {
    const { container } = render(
      <Grid
        columns={columns}
        rows={[["Alice", "active"]]}
        height={300}
        totalRecords={3}
      />,
    );

    expect(screen.getByText("Alice")).toBeInTheDocument();
    // Two unloaded rows × two columns = four skeleton cells
    expect(container.querySelectorAll(".h-4.w-full").length).toBe(4);
  });

  it("renders trailingHeader in the table header", () => {
    render(
      <Grid
        columns={columns}
        rows={rows}
        height={300}
        trailingHeader={<button type="button">Add column</button>}
      />,
    );
    expect(
      screen.getByRole("button", { name: "Add column" }),
    ).toBeInTheDocument();
  });

  it("renders an AI icon for columns marked isAI", () => {
    render(
      <Grid
        columns={[{ name: "Summary", isAI: true }]}
        rows={[["hello"]]}
        height={300}
      />,
    );
    expect(screen.getByText("auto_awesome")).toBeInTheDocument();
    expect(screen.getByText("Summary")).toBeInTheDocument();
  });
});
