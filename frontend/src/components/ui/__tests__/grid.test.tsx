import { render } from "@testing-library/react";
import Grid from "../grid";

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
});
