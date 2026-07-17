import { render } from "@testing-library/react";
import { AspectRatio } from "../aspect-ratio";

describe("AspectRatio", () => {
  it("renders without crashing", () => {
    render(<AspectRatio ratio={16 / 9} />);
  });
});
