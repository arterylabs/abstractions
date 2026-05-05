import { render } from "@testing-library/react";
import { HoverCard, HoverCardTrigger, HoverCardContent } from "../hover-card";

describe("HoverCard", () => {
  it("renders root without crashing", () => {
    render(
      <HoverCard>
        <HoverCardTrigger>Hover me</HoverCardTrigger>
        <HoverCardContent>Content</HoverCardContent>
      </HoverCard>,
    );
  });
});
