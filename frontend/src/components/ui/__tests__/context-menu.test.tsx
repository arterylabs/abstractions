import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  ContextMenu,
  ContextMenuTrigger,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuCheckboxItem,
  ContextMenuRadioGroup,
  ContextMenuRadioItem,
  ContextMenuLabel,
  ContextMenuSeparator,
  ContextMenuShortcut,
  ContextMenuGroup,
} from "../context-menu";

function renderContextMenu(children: React.ReactNode) {
  return render(
    <ContextMenu>
      <ContextMenuTrigger>Trigger</ContextMenuTrigger>
      <ContextMenuContent>{children}</ContextMenuContent>
    </ContextMenu>,
  );
}

function openContextMenu() {
  fireEvent.contextMenu(screen.getByText("Trigger"));
}

describe("ContextMenu", () => {
  it("renders without crashing", () => {
    render(
      <ContextMenu>
        <ContextMenuTrigger>Right click me</ContextMenuTrigger>
      </ContextMenu>,
    );
  });

  it("renders trigger text", () => {
    render(
      <ContextMenu>
        <ContextMenuTrigger>Right click me</ContextMenuTrigger>
      </ContextMenu>,
    );
    expect(screen.getByText("Right click me")).toBeInTheDocument();
  });

  it("renders content items after right-click", async () => {
    renderContextMenu(
      <>
        <ContextMenuItem>Open</ContextMenuItem>
        <ContextMenuItem>Copy</ContextMenuItem>
      </>,
    );
    openContextMenu();
    expect(await screen.findByText("Open")).toBeInTheDocument();
    expect(screen.getByText("Copy")).toBeInTheDocument();
  });

  it("renders ContextMenuLabel and ContextMenuSeparator after right-click", async () => {
    renderContextMenu(
      <>
        <ContextMenuLabel>Actions</ContextMenuLabel>
        <ContextMenuSeparator />
        <ContextMenuGroup>
          <ContextMenuItem>Edit</ContextMenuItem>
        </ContextMenuGroup>
      </>,
    );
    openContextMenu();
    expect(await screen.findByText("Actions")).toBeInTheDocument();
    expect(screen.getByText("Edit")).toBeInTheDocument();
  });

  it("renders ContextMenuCheckboxItem with checked state", async () => {
    renderContextMenu(
      <ContextMenuCheckboxItem checked={true}>
        Show labels
      </ContextMenuCheckboxItem>,
    );
    openContextMenu();
    expect(await screen.findByText("Show labels")).toBeInTheDocument();
  });

  it("fires onCheckedChange when checkbox item is clicked", async () => {
    const user = userEvent.setup();
    const onCheckedChange = jest.fn();
    renderContextMenu(
      <ContextMenuCheckboxItem
        checked={false}
        onCheckedChange={onCheckedChange}
      >
        Toggle me
      </ContextMenuCheckboxItem>,
    );
    openContextMenu();
    await user.click(await screen.findByText("Toggle me"));
    expect(onCheckedChange).toHaveBeenCalledTimes(1);
  });

  it("renders ContextMenuRadioGroup with radio items", async () => {
    renderContextMenu(
      <ContextMenuRadioGroup value="light">
        <ContextMenuRadioItem value="light">Light</ContextMenuRadioItem>
        <ContextMenuRadioItem value="dark">Dark</ContextMenuRadioItem>
      </ContextMenuRadioGroup>,
    );
    openContextMenu();
    expect(await screen.findByText("Light")).toBeInTheDocument();
    expect(screen.getByText("Dark")).toBeInTheDocument();
  });

  it("fires onValueChange when a radio item is selected", async () => {
    const user = userEvent.setup();
    const onValueChange = jest.fn();
    renderContextMenu(
      <ContextMenuRadioGroup value="light" onValueChange={onValueChange}>
        <ContextMenuRadioItem value="light">Light</ContextMenuRadioItem>
        <ContextMenuRadioItem value="dark">Dark</ContextMenuRadioItem>
      </ContextMenuRadioGroup>,
    );
    openContextMenu();
    await user.click(await screen.findByText("Dark"));
    expect(onValueChange).toHaveBeenCalledWith("dark");
  });

  it("renders ContextMenuShortcut inside an item", async () => {
    renderContextMenu(
      <ContextMenuItem>
        Rename
        <ContextMenuShortcut>F2</ContextMenuShortcut>
      </ContextMenuItem>,
    );
    openContextMenu();
    expect(await screen.findByText("F2")).toBeInTheDocument();
  });

  it("fires onClick when a menu item is clicked", async () => {
    const user = userEvent.setup();
    const handleClick = jest.fn();
    renderContextMenu(
      <ContextMenuItem onClick={handleClick}>Delete</ContextMenuItem>,
    );
    openContextMenu();
    await user.click(await screen.findByText("Delete"));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});
