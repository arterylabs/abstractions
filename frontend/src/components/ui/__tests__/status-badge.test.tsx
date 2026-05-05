import { render, screen } from "@testing-library/react";
import { StatusBadge } from "../status-badge";

describe("StatusBadge", () => {
  it("should render 'Pending' for status 'pending'", () => {
    render(<StatusBadge status="pending" />);
    expect(screen.getByText("Pending")).toBeInTheDocument();
  });

  it("should render 'Accepted' for status 'accepted'", () => {
    render(<StatusBadge status="accepted" />);
    expect(screen.getByText("Accepted")).toBeInTheDocument();
  });

  it("should use custom label when provided", () => {
    render(<StatusBadge status="pending" label="Awaiting Review" />);
    expect(screen.getByText("Awaiting Review")).toBeInTheDocument();
  });

  it("should render the raw status string when status is unknown", () => {
    render(<StatusBadge status="unknown_status" />);
    expect(screen.getByText("unknown_status")).toBeInTheDocument();
  });

  it("should be case-insensitive (ACTIVE → Active)", () => {
    render(<StatusBadge status="ACTIVE" />);
    expect(screen.getByText("Active")).toBeInTheDocument();
  });
});
