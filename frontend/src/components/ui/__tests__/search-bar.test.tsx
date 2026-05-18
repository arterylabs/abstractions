import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { SearchBar } from "../search-bar";

describe("SearchBar", () => {
  it("renders an input element", () => {
    render(<SearchBar />);
    expect(screen.getByRole("textbox")).toBeInTheDocument();
  });

  it("renders with a custom placeholder", () => {
    render(<SearchBar placeholder="Find anything..." />);
    expect(screen.getByPlaceholderText("Find anything...")).toBeInTheDocument();
  });

  it("calls onChange when the user types", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<SearchBar onChange={onChange} />);
    await user.type(screen.getByRole("textbox"), "a");
    expect(onChange).toHaveBeenCalledWith("a");
  });
});
