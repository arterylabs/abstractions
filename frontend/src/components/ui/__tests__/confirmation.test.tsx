import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ConfirmationDialog } from "../confirmation";

describe("ConfirmationDialog", () => {
  const defaultProps = {
    title: "Confirm Action",
    content: "Are you sure you want to proceed?",
    okCallback: jest.fn().mockResolvedValue(undefined),
    onClose: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("should render title and content when open=true", () => {
    render(<ConfirmationDialog {...defaultProps} open={true} />);
    expect(screen.getByText("Confirm Action")).toBeInTheDocument();
    expect(
      screen.getByText("Are you sure you want to proceed?"),
    ).toBeInTheDocument();
  });

  it("should not render content when open=false", () => {
    render(<ConfirmationDialog {...defaultProps} open={false} />);
    expect(screen.queryByText("Confirm Action")).not.toBeInTheDocument();
  });

  it("should call onClose when cancel button is clicked", async () => {
    const user = userEvent.setup();
    const onClose = jest.fn();
    render(
      <ConfirmationDialog {...defaultProps} open={true} onClose={onClose} />,
    );
    await user.click(screen.getByRole("button", { name: "Cancel" }));
    expect(onClose).toHaveBeenCalled();
  });

  it("should call okCallback and then onClose when confirm button is clicked", async () => {
    const user = userEvent.setup();
    const okCallback = jest.fn().mockResolvedValue(undefined);
    const onClose = jest.fn();
    render(
      <ConfirmationDialog
        {...defaultProps}
        open={true}
        okCallback={okCallback}
        onClose={onClose}
      />,
    );
    await user.click(screen.getByRole("button", { name: "Confirm" }));
    await waitFor(() => {
      expect(okCallback).toHaveBeenCalledTimes(1);
      expect(onClose).toHaveBeenCalled();
    });
  });

  it("should show 'Processing...' text while okCallback is pending", async () => {
    const user = userEvent.setup();
    let resolveCallback!: () => void;
    const okCallback = jest.fn().mockReturnValue(
      new Promise<void>((resolve) => {
        resolveCallback = resolve;
      }),
    );
    render(
      <ConfirmationDialog
        {...defaultProps}
        open={true}
        okCallback={okCallback}
        onClose={jest.fn()}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Confirm" }));
    await waitFor(() => {
      expect(screen.getByText("Processing...")).toBeInTheDocument();
    });

    // Resolve and flush the resulting setIsLoading(false) state update
    await act(async () => {
      resolveCallback();
    });
  });

  it("should render confirm button without crashing with destructive variant", () => {
    render(
      <ConfirmationDialog
        {...defaultProps}
        open={true}
        variant="destructive"
      />,
    );
    expect(screen.getByRole("button", { name: "Confirm" })).toBeInTheDocument();
  });

  it("should render custom okText and cancelText", () => {
    render(
      <ConfirmationDialog
        {...defaultProps}
        open={true}
        okText="Delete"
        cancelText="Keep"
      />,
    );
    expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Keep" })).toBeInTheDocument();
  });
});
