import { render } from "@testing-library/react";
import { ContextMenu, ContextMenuTrigger } from "../context-menu";

describe("ContextMenu", () => {
  it("renders without crashing", () => {
    render(
      <ContextMenu>
        <ContextMenuTrigger>Right click me</ContextMenuTrigger>
      </ContextMenu>,
    );
  });
});
