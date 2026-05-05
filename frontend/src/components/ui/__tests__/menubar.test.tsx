import { render } from "@testing-library/react";
import { Menubar, MenubarMenu, MenubarTrigger } from "../menubar";

describe("Menubar", () => {
  it("renders without crashing", () => {
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>File</MenubarTrigger>
        </MenubarMenu>
      </Menubar>,
    );
  });
});
