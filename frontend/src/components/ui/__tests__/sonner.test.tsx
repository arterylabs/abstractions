import { render } from "@testing-library/react";
import { Toaster } from "../sonner";

describe("Sonner Toaster", () => {
  it("renders without crashing", () => {
    render(<Toaster />);
  });
});
