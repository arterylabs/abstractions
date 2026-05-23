import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
  MenubarContent,
  MenubarItem,
  MenubarCheckboxItem,
  MenubarRadioGroup,
  MenubarRadioItem,
  MenubarLabel,
  MenubarSeparator,
  MenubarShortcut,
  MenubarGroup,
} from "../menubar";

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

  it("renders menu trigger text", () => {
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>Edit</MenubarTrigger>
        </MenubarMenu>
      </Menubar>,
    );
    expect(screen.getByText("Edit")).toBeInTheDocument();
  });

  it("renders MenubarContent with items after clicking trigger", async () => {
    const user = userEvent.setup();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>File</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>New File</MenubarItem>
            <MenubarItem>Open...</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("File"));
    expect(await screen.findByText("New File")).toBeInTheDocument();
    expect(screen.getByText("Open...")).toBeInTheDocument();
  });

  it("fires onClick when a menu item is clicked", async () => {
    const user = userEvent.setup();
    const handleClick = jest.fn();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>File</MenubarTrigger>
          <MenubarContent>
            <MenubarItem onClick={handleClick}>Save</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("File"));
    await user.click(await screen.findByText("Save"));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("renders MenubarLabel and MenubarSeparator after opening", async () => {
    const user = userEvent.setup();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>View</MenubarTrigger>
          <MenubarContent>
            <MenubarLabel>Appearance</MenubarLabel>
            <MenubarSeparator />
            <MenubarGroup>
              <MenubarItem>Zoom In</MenubarItem>
            </MenubarGroup>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("View"));
    expect(await screen.findByText("Appearance")).toBeInTheDocument();
    expect(screen.getByText("Zoom In")).toBeInTheDocument();
  });

  it("renders MenubarCheckboxItem after opening", async () => {
    const user = userEvent.setup();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>View</MenubarTrigger>
          <MenubarContent>
            <MenubarCheckboxItem checked={true}>
              Show Toolbar
            </MenubarCheckboxItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("View"));
    expect(await screen.findByText("Show Toolbar")).toBeInTheDocument();
  });

  it("fires onCheckedChange when checkbox item is clicked", async () => {
    const user = userEvent.setup();
    const onCheckedChange = jest.fn();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>View</MenubarTrigger>
          <MenubarContent>
            <MenubarCheckboxItem
              checked={false}
              onCheckedChange={onCheckedChange}
            >
              Show Sidebar
            </MenubarCheckboxItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("View"));
    await user.click(await screen.findByText("Show Sidebar"));
    expect(onCheckedChange).toHaveBeenCalledTimes(1);
  });

  it("renders MenubarRadioGroup with items after opening", async () => {
    const user = userEvent.setup();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>Theme</MenubarTrigger>
          <MenubarContent>
            <MenubarRadioGroup value="light">
              <MenubarRadioItem value="light">Light</MenubarRadioItem>
              <MenubarRadioItem value="dark">Dark</MenubarRadioItem>
            </MenubarRadioGroup>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("Theme"));
    expect(await screen.findByText("Light")).toBeInTheDocument();
    expect(screen.getByText("Dark")).toBeInTheDocument();
  });

  it("fires onValueChange when a radio item is selected", async () => {
    const user = userEvent.setup();
    const onValueChange = jest.fn();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>Theme</MenubarTrigger>
          <MenubarContent>
            <MenubarRadioGroup value="light" onValueChange={onValueChange}>
              <MenubarRadioItem value="light">Light</MenubarRadioItem>
              <MenubarRadioItem value="dark">Dark</MenubarRadioItem>
            </MenubarRadioGroup>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("Theme"));
    await user.click(await screen.findByText("Dark"));
    expect(onValueChange).toHaveBeenCalledWith("dark");
  });

  it("renders MenubarShortcut inside an item after opening", async () => {
    const user = userEvent.setup();
    render(
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>File</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>
              Save
              <MenubarShortcut>⌘S</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>,
    );
    await user.click(screen.getByText("File"));
    expect(await screen.findByText("⌘S")).toBeInTheDocument();
  });
});
