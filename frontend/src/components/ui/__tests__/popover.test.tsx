import { render } from "@testing-library/react";
import { Popover, PopoverTrigger, PopoverContent } from "../popover";

describe("Popover", () => {
  it("renders root without crashing (closed state)", () => {
    render(
      <Popover>
        <PopoverTrigger>Open</PopoverTrigger>
      </Popover>,
    );
  });

  it("renders content without crashing when open", () => {
    render(
      <Popover open>
        <PopoverTrigger>Open</PopoverTrigger>
        <PopoverContent>Content</PopoverContent>
      </Popover>,
    );
  });
});
