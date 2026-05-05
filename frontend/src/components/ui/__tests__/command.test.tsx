import { render } from "@testing-library/react";
import {
  Command,
  CommandInput,
  CommandList,
  CommandItem,
  CommandEmpty,
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
});
