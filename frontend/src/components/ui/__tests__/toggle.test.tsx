import { render } from "@testing-library/react";
import { Toggle } from "../toggle";

describe("Toggle", () => {
  it("renders without crashing", () => {
    render(<Toggle>Bold</Toggle>);
  });
});
