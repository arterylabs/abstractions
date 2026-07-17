import { render, screen, act } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { SearchInput } from "../search-input";

describe("SearchInput", () => {
  afterEach(() => {
    jest.useRealTimers();
  });

  it("should render input with placeholder", () => {
    render(<SearchInput placeholder="Search items..." />);
    expect(screen.getByPlaceholderText("Search items...")).toBeInTheDocument();
  });

  it("should call onChange when typing", async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<SearchInput onChange={onChange} />);
    await user.type(screen.getByRole("textbox"), "a");
    expect(onChange).toHaveBeenCalledWith("a");
  });

  it("should call onSearch after debounceMs with the typed value", async () => {
    jest.useFakeTimers();
    const user = userEvent.setup({ delay: null });
    const onSearch = jest.fn();
    render(<SearchInput onSearch={onSearch} debounceMs={300} />);

    await user.type(screen.getByRole("textbox"), "hello");
    expect(onSearch).not.toHaveBeenCalled();

    act(() => {
      jest.runAllTimers();
    });

    expect(onSearch).toHaveBeenLastCalledWith("hello");
  });

  it("should render provided value in controlled mode", () => {
    render(<SearchInput value="controlled text" onChange={jest.fn()} />);
    expect(screen.getByDisplayValue("controlled text")).toBeInTheDocument();
  });

  it("should update displayed value in uncontrolled mode when typing", async () => {
    const user = userEvent.setup();
    render(<SearchInput />);
    await user.type(screen.getByRole("textbox"), "typed text");
    expect(screen.getByDisplayValue("typed text")).toBeInTheDocument();
  });
});
