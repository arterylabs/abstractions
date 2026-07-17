import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Toaster, Notification } from "../sonner";

describe("Sonner Toaster", () => {
  it("renders without crashing", () => {
    render(<Toaster />);
  });
});

describe("Notification", () => {
  it("renders success type with title", () => {
    render(<Notification type="success" params={{ title: "Saved!" }} />);
    expect(screen.getByText("Saved!")).toBeInTheDocument();
  });

  it("renders error type with title", () => {
    render(
      <Notification type="error" params={{ title: "Something went wrong" }} />,
    );
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders info type with title", () => {
    render(<Notification type="info" params={{ title: "FYI" }} />);
    expect(screen.getByText("FYI")).toBeInTheDocument();
  });

  it("renders warning type with title", () => {
    render(<Notification type="warning" params={{ title: "Heads up" }} />);
    expect(screen.getByText("Heads up")).toBeInTheDocument();
  });

  it("renders optional message when provided", () => {
    render(
      <Notification
        type="success"
        params={{ title: "Done", message: "Your changes have been saved." }}
      />,
    );
    expect(screen.getByText("Done")).toBeInTheDocument();
    expect(
      screen.getByText("Your changes have been saved."),
    ).toBeInTheDocument();
  });

  it("does not render message element when message is omitted", () => {
    render(<Notification type="info" params={{ title: "Title only" }} />);
    expect(screen.queryByText(/saved/i)).not.toBeInTheDocument();
  });

  it("renders action buttons and fires onClick", async () => {
    const user = userEvent.setup();
    const handleRetry = jest.fn();
    render(
      <Notification
        type="error"
        params={{
          title: "Failed",
          actions: [{ label: "Retry", onClick: handleRetry }],
        }}
      />,
    );
    await user.click(screen.getByRole("button", { name: "Retry" }));
    expect(handleRetry).toHaveBeenCalledTimes(1);
  });

  it("renders multiple action buttons", () => {
    render(
      <Notification
        type="info"
        params={{
          title: "Update available",
          actions: [
            { label: "Update now", onClick: jest.fn() },
            { label: "Later", onClick: jest.fn() },
          ],
        }}
      />,
    );
    expect(
      screen.getByRole("button", { name: "Update now" }),
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Later" })).toBeInTheDocument();
  });

  it("accepts a plain string as params", () => {
    render(<Notification type="success" params="Quick message" />);
    expect(screen.getByText("Quick message")).toBeInTheDocument();
  });
});
