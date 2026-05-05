import { render } from "@testing-library/react";
import { ScrollArea, ScrollBar } from "../scroll-area";

describe("ScrollArea", () => {
  it("renders without crashing", () => {
    render(
      <ScrollArea>
        <div>Scrollable content</div>
      </ScrollArea>,
    );
  });

  it("renders ScrollBar inside ScrollArea without crashing", () => {
    render(
      <ScrollArea>
        <div>content</div>
        <ScrollBar />
      </ScrollArea>,
    );
  });
});
