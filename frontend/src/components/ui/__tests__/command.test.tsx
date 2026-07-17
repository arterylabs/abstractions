import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  Command,
  CommandDialog,
  CommandInput,
  CommandList,
  CommandItem,
  CommandEmpty,
  CommandGroup,
  CommandSeparator,
  CommandShortcut,
} from "../command";

describe("Command", () => {
  it("renders without crashing", () => {
    render(
      <Command>
        <CommandInput placeholder="Search..." />
        <CommandList>
          <CommandEmpty>No results.</CommandEmpty>
          <CommandItem>Item one</CommandItem>
        </CommandList>
      </Command>,
    );
  });

  it("renders CommandGroup with a heading and items", () => {
    render(
      <Command>
        <CommandList>
          <CommandGroup heading="Suggestions">
            <CommandItem>Calendar</CommandItem>
            <CommandItem>Search</CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>,
    );
    expect(screen.getByText("Suggestions")).toBeInTheDocument();
    expect(screen.getByText("Calendar")).toBeInTheDocument();
    expect(screen.getByText("Search")).toBeInTheDocument();
  });

  it("renders CommandSeparator between groups", () => {
    const { container } = render(
      <Command>
        <CommandList>
          <CommandGroup heading="Group A">
            <CommandItem>Item A</CommandItem>
          </CommandGroup>
          <CommandSeparator />
          <CommandGroup heading="Group B">
            <CommandItem>Item B</CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>,
    );
    const separator = container.querySelector(
      '[data-slot="command-separator"]',
    );
    expect(separator).toBeInTheDocument();
  });

  it("renders CommandShortcut inside an item", () => {
    render(
      <Command>
        <CommandList>
          <CommandItem>
            Open file
            <CommandShortcut>⌘O</CommandShortcut>
          </CommandItem>
        </CommandList>
      </Command>,
    );
    expect(screen.getByText("⌘O")).toBeInTheDocument();
  });

  it("shows CommandEmpty when no items match the search", async () => {
    const user = userEvent.setup();
    render(
      <Command>
        <CommandInput placeholder="Search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandItem>Apple</CommandItem>
        </CommandList>
      </Command>,
    );
    await user.type(screen.getByPlaceholderText("Search..."), "zzz");
    expect(screen.getByText("No results found.")).toBeInTheDocument();
  });

  it("renders CommandDialog when open", () => {
    render(
      <CommandDialog open title="Command Palette" description="Run a command">
        <CommandList>
          <CommandItem>New file</CommandItem>
        </CommandList>
      </CommandDialog>,
    );
    expect(screen.getByText("New file")).toBeInTheDocument();
  });

  it("CommandDialog renders custom title and description accessibly", () => {
    render(
      <CommandDialog
        open
        title="My Palette"
        description="Type to search actions"
      >
        <CommandList>
          <CommandEmpty>Nothing here.</CommandEmpty>
        </CommandList>
      </CommandDialog>,
    );
    expect(screen.getByText("My Palette")).toBeInTheDocument();
    expect(screen.getByText("Type to search actions")).toBeInTheDocument();
  });
});
