import React from "react";
import { render, screen } from "@testing-library/react";
import {
  ChartContainer,
  ChartTooltipContent,
  ChartLegendContent,
} from "../chart";
import type { ChartConfig } from "../chart";
import { LineChart, Line } from "recharts";

const config: ChartConfig = {
  value: { label: "Value", color: "#3b82f6" },
};

// Provides ChartContext by rendering the component inside ChartContainer.
// ChartContainer passes children to ResponsiveContainer which renders them
// in jsdom (at width=0/height=0 with ResizeObserver mocked in setupTests).
function renderInChart(
  ui: React.ReactElement,
  chartConfig: ChartConfig = config,
) {
  return render(
    <ChartContainer config={chartConfig}>
      {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
      {ui as React.ReactElement<any, any>}
    </ChartContainer>,
  );
}

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

  it("injects CSS color variable when config has a color entry", () => {
    const { container } = render(
      <ChartContainer config={config}>
        <LineChart data={[]}>
          <Line dataKey="value" />
        </LineChart>
      </ChartContainer>,
    );
    const styleEl = container.querySelector("style");
    expect(styleEl).not.toBeNull();
    expect(styleEl!.textContent).toContain("--color-value");
    expect(styleEl!.textContent).toContain("#3b82f6");
  });

  it("omits style element when config has no color or theme", () => {
    const noColorConfig: ChartConfig = { value: { label: "No Color" } };
    const { container } = render(
      <ChartContainer config={noColorConfig}>
        <LineChart data={[]}>
          <Line dataKey="value" />
        </LineChart>
      </ChartContainer>,
    );
    expect(container.querySelector("style")).toBeNull();
  });

  it("ChartTooltipContent renders payload value when active", () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any[] = [
      {
        dataKey: "value",
        name: "value",
        value: 42,
        color: "#3b82f6",
        payload: { fill: "#3b82f6" },
      },
    ];
    renderInChart(<ChartTooltipContent active={true} payload={payload} />);
    expect(screen.getByText("42")).toBeInTheDocument();
  });

  it("ChartTooltipContent resolves label from config when label matches a config key", () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any[] = [
      {
        dataKey: "value",
        name: "value",
        value: 10,
        payload: { fill: "#3b82f6" },
      },
    ];
    renderInChart(
      <ChartTooltipContent active={true} payload={payload} label="value" />,
    );
    // "Value" appears as the resolved tooltip label (from config.value.label)
    expect(screen.getAllByText("Value").length).toBeGreaterThan(0);
  });

  it("ChartTooltipContent renders nothing when not active", () => {
    renderInChart(<ChartTooltipContent active={false} payload={[]} />);
    // Inactive tooltip returns null — no value or label text rendered
    expect(screen.queryByText("42")).not.toBeInTheDocument();
  });

  it("ChartLegendContent renders config label for each payload item", () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const legendPayload: any[] = [
      { value: "value", dataKey: "value", color: "#3b82f6" },
    ];
    renderInChart(<ChartLegendContent payload={legendPayload} />);
    expect(screen.getByText("Value")).toBeInTheDocument();
  });

  it("ChartLegendContent renders nothing with an empty payload", () => {
    const { container } = renderInChart(<ChartLegendContent payload={[]} />);
    // Null render — only the ChartContainer wrapper div remains
    expect(
      container.querySelector(
        '[class*="flex items-center justify-center gap-4"]',
      ),
    ).toBeNull();
  });
});
