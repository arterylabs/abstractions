import { render, screen } from "@testing-library/react";
import { DataTable } from "../data-table";
import type { DataTableColumn } from "../data-table";

const columns: DataTableColumn[] = [
  { key: "name", label: "Name", order: 0 },
  { key: "status", label: "Status", order: 1 },
];

const rows = [
  { name: "Alice", status: "active" },
  { name: "Bob", status: "inactive" },
];

describe("DataTable", () => {
  it("renders without crashing", () => {
    render(<DataTable columns={columns} data={rows} />);
  });

  it("renders column headers from label", () => {
    render(<DataTable columns={columns} data={rows} />);
    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Status")).toBeInTheDocument();
  });

  it("renders columns sorted by order regardless of array order", () => {
    const reversedColumns: DataTableColumn[] = [
      { key: "status", label: "Status", order: 1 },
      { key: "name", label: "Name", order: 0 },
    ];
    render(<DataTable columns={reversedColumns} data={rows} />);
    const headers = screen.getAllByRole("columnheader");
    expect(headers[0]).toHaveTextContent("Name");
    expect(headers[1]).toHaveTextContent("Status");
  });

  it("shows spinner when loading=true", () => {
    const { container } = render(
      <DataTable columns={columns} data={rows} loading={true} />,
    );
    expect(container.querySelector("svg")).toBeInTheDocument();
  });

  it("shows default empty message when data is empty", () => {
    render(<DataTable columns={columns} data={[]} />);
    expect(screen.getByText("No data available")).toBeInTheDocument();
  });

  it("shows custom emptyMessage when data is empty", () => {
    render(
      <DataTable columns={columns} data={[]} emptyMessage="Nothing here yet" />,
    );
    expect(screen.getByText("Nothing here yet")).toBeInTheDocument();
  });

  it("renders null cell value as em dash", () => {
    const nullRows = [{ name: null, status: "active" }];
    render(<DataTable columns={columns} data={nullRows} />);
    expect(screen.getByText("—")).toBeInTheDocument();
  });

  it("renders undefined cell value as em dash", () => {
    const undefinedRows = [{ name: undefined, status: "active" }];
    render(<DataTable columns={columns} data={undefinedRows} />);
    expect(screen.getByText("—")).toBeInTheDocument();
  });

  it("renders boolean true as 'Yes'", () => {
    const boolColumns: DataTableColumn[] = [
      { key: "active", label: "Active", order: 0 },
    ];
    render(<DataTable columns={boolColumns} data={[{ active: true }]} />);
    expect(screen.getByText("Yes")).toBeInTheDocument();
  });

  it("renders boolean false as 'No'", () => {
    const boolColumns: DataTableColumn[] = [
      { key: "active", label: "Active", order: 0 },
    ];
    render(<DataTable columns={boolColumns} data={[{ active: false }]} />);
    expect(screen.getByText("No")).toBeInTheDocument();
  });

  it("uses custom render function for a column", () => {
    const customColumns: DataTableColumn[] = [
      {
        key: "name",
        label: "Name",
        order: 0,
        render: (value) => <strong>{String(value)}-custom</strong>,
      },
    ];
    render(<DataTable columns={customColumns} data={[{ name: "Alice" }]} />);
    expect(screen.getByText("Alice-custom")).toBeInTheDocument();
  });

  it("uses custom headerRender function for a column", () => {
    const customColumns: DataTableColumn[] = [
      {
        key: "name",
        label: "Name",
        order: 0,
        headerRender: () => <span>Custom Header</span>,
      },
    ];
    render(<DataTable columns={customColumns} data={rows} />);
    expect(screen.getByText("Custom Header")).toBeInTheDocument();
  });

  it("applies rowClassName to each row", () => {
    const { container } = render(
      <DataTable
        columns={columns}
        data={rows}
        rowClassName={(row) =>
          (row as { status: string }).status === "active" ? "row-active" : ""
        }
      />,
    );
    const activeRow = container.querySelector("tr.row-active");
    expect(activeRow).toBeInTheDocument();
  });
});
