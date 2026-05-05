import { render } from "@testing-library/react";
import { ChartContainer } from "../chart";
import type { ChartConfig } from "../chart";
import { LineChart, Line } from "recharts";

const config: ChartConfig = {
  value: { label: "Value", color: "#3b82f6" },
};

describe("Chart", () => {
  beforeEach(() => {
    // jsdom has no layout engine so recharts ResponsiveContainer always sees
    // width=0/height=0 and emits a console.warn — suppress it in tests.
    jest.spyOn(console, "warn").mockImplementation((msg, ...rest) => {
      if (typeof msg === "string" && msg.includes("width(0)")) return;
      console.warn(msg, ...rest);
    });
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it("renders ChartContainer without crashing", () => {
    render(
      <ChartContainer config={config}>
        <LineChart data={[{ value: 1 }, { value: 2 }]}>
          <Line dataKey="value" />
        </LineChart>
      </ChartContainer>,
    );
  });
});
