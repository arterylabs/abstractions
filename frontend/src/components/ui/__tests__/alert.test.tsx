import { render, screen } from "@testing-library/react";
import { Alert, AlertTitle, AlertDescription } from "../alert";

describe("Alert", () => {
  it("renders without crashing", () => {
    render(
      <Alert>
        <AlertTitle>Title</AlertTitle>
        <AlertDescription>Description</AlertDescription>
      </Alert>,
    );
  });

  it("renders with role='alert'", () => {
    render(<Alert>Alert content</Alert>);
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });

  it("renders AlertTitle text", () => {
    render(
      <Alert>
        <AlertTitle>Something went wrong</AlertTitle>
      </Alert>,
    );
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders AlertDescription content", () => {
    render(
      <Alert>
        <AlertDescription>Please try again later.</AlertDescription>
      </Alert>,
    );
    expect(screen.getByText("Please try again later.")).toBeInTheDocument();
  });

  it("renders destructive variant without crashing", () => {
    render(
      <Alert variant="destructive">
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>An error occurred.</AlertDescription>
      </Alert>,
    );
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });

  it("renders default variant without crashing", () => {
    render(
      <Alert variant="default">
        <AlertTitle>Info</AlertTitle>
      </Alert>,
    );
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });
});
